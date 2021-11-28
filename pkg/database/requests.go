package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/impulse-http/local-backend/pkg/models"
	"log"
	"net/http"
	"time"
)

func insertHeadersValues(ctx context.Context, db *sql.DB, requestHistoryId, requestId int64, headers http.Header, isRequest bool) error {
	query := `
			INSERT INTO headers (key, request_history_id, is_request)
			VALUES ($1, $2, $3);
	`
	writeId := requestHistoryId
	if requestId > 0 {
		query = `
			INSERT INTO headers (key, request_id, is_request)
			VALUES ($1, $2, $3);
		`
		writeId = requestId
	}
	for name, values := range headers {
		r, err := db.ExecContext(ctx,
			query,
			name,
			writeId,
			boolToInt(isRequest),
		)

		if err != nil {
			log.Println("Error while inserting header" + name)
			return err
		}

		headerId, err := r.LastInsertId()

		if err != nil {
			log.Println("Error getting row id for header " + name)
			return err
		}

		for _, value := range values {
			r, err = db.ExecContext(ctx,
				`
			INSERT INTO headers_values (header_id, header_value)
			VALUES ($1, $2);
			`,
				headerId,
				value,
			)

			if err != nil {
				log.Println("Error inserting header value for key = " + name + " value = " + value)
				return err
			}
		}
	}

	return nil
}

func deleteHeadersValues(ctx context.Context, db *sql.DB, requestHistoryId, requestId int64) error {
	query := `DELETE FROM headers WHERE request_history_id = $1`
	deleteId := requestHistoryId
	if requestId > 0 {
		query = `DELETE FROM headers WHERE request_id = $1`
		deleteId = requestId
	}
	_, err := db.ExecContext(ctx, query, deleteId)
	return err
}

func updateRequestHeadersValues(ctx context.Context, db *sql.DB, requestId int64, reqHeaders http.Header) error {
	if err := deleteHeadersValues(ctx, db, 0, requestId); err != nil {
		return err
	}
	if err := insertHeadersValues(ctx, db, 0, requestId, reqHeaders, true); err != nil {
		return err
	}
	return nil
}

func getRequestHeaders(ctx context.Context, db *sql.DB, requestHistoryId, requestId int64) (http.Header, http.Header, error) {
	query := `
		SELECT h.key, hv.header_value, h.is_request
		FROM headers h JOIN headers_values hv ON h.id = hv.header_id
		WHERE h.request_history_id = $1
	`
	getId := requestHistoryId

	if requestId > 0 {
		query = `
			SELECT h.key, hv.header_value, h.is_request
			FROM headers h JOIN headers_values hv ON h.id = hv.header_id
			WHERE request_id = $1
		`
		getId = requestId
	}

	r, err := db.QueryContext(ctx, query, getId)

	reqHeaders := make(http.Header)
	resHeaders := make(http.Header)
	key := ""
	value := ""
	isRequest := 0

	for r.Next() {
		if err := r.Scan(&key, &value, &isRequest); err != nil {
			return nil, nil, err
		}

		if intToBool(isRequest) {
			reqHeaders[key] = append(reqHeaders[key], value)
		} else {
			resHeaders[key] = append(resHeaders[key], value)
		}
	}

	return reqHeaders, resHeaders, err
}

func (d *Database) CreateHistoryEntry(ctx context.Context, req *models.RequestType, res *models.ResponseType) (int64, error) {
	db := d.db
	r, err := db.ExecContext(ctx,
		`
		INSERT INTO requests_history (request_body, response_body, user_id, created_at, method)
		VALUES ($1, $2, 1, $3, $4)
		`,
		req.Body,
		res.Body,
		time.Now().Unix(),
		req.Method,
	)
	if err != nil {
		log.Println("Error running query: " + err.Error())
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		log.Println("Error fetching last row id: " + err.Error())
		return 0, err
	}

	err = insertHeadersValues(ctx, db, id, 0, req.Headers, true)
	if err != nil {
		return id, err
	}

	err = insertHeadersValues(ctx, db, id, 0, res.Headers, false)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (d *Database) GetHistory(ctx context.Context) ([]models.RequestHistoryEntry, error) {
	type rawHistory struct {
		Id           int64
		Method       string
		Body         string
		ResponseBody string
		CreatedAt    time.Time
		IsRequest    int
		HeaderKey    string
		HeaderValue  string
	}

	newHistoryEntry := func() models.RequestHistoryEntry {
		entry := models.RequestHistoryEntry{}
		entry.Request.Headers = make(http.Header)
		entry.Response.Headers = make(http.Header)
		return entry
	}

	db := d.db
	rows, err := db.QueryContext(ctx, `
			SELECT rh.id, rh.method, rh.request_body, rh.response_body, rh.created_at,
			       h.is_request, h.key, hv.header_value
			FROM requests_history rh
			    JOIN headers h on rh.id = h.request_history_id
				JOIN headers_values hv ON hv.header_id = h.id
			ORDER BY rh.id, h.id;
	`)

	if err != nil {
		log.Println("Request history query failed" + fmt.Sprint(err))
		return nil, err
	}

	res := make([]models.RequestHistoryEntry, 0)
	cur := newHistoryEntry()

	for rows.Next() {
		row := rawHistory{}
		err = rows.Scan(
			&row.Id,
			&row.Method,
			&row.Body,
			&row.ResponseBody,
			&row.CreatedAt,
			&row.IsRequest,
			&row.HeaderKey,
			&row.HeaderValue,
		)

		if err != nil {
			return nil, err
		}

		if cur.Id == 0 || row.Id == cur.Id {
			if cur.Id == 0 {
				cur.Id = row.Id
				cur.CreatedAt = row.CreatedAt
				cur.Request.Body = row.Body
				cur.Request.Method = row.Method
				cur.Response.Body = row.ResponseBody
			}

			if intToBool(row.IsRequest) {
				cur.Request.Headers[row.HeaderKey] = append(cur.Request.Headers[row.HeaderKey], row.HeaderValue)
			} else {
				cur.Response.Headers[row.HeaderKey] = append(cur.Response.Headers[row.HeaderKey], row.HeaderValue)
			}
		} else {
			res = append(res, cur)
			cur = newHistoryEntry()
		}
	}

	return res, nil
}

func (d *Database) CreateRequest(ctx context.Context, request *models.Request) (int64, error) {
	query := `INSERT INTO requests(name, request_body, user_id, created_at, method) VALUES ($1, $2, 1, $3, $4)`
	if request.CollectionId > 0 {
		query = `
			INSERT INTO requests(
				 name,
				 request_body,
				 user_id,
				 created_at,
				 method,
				 collection_id
			)
			VALUES ($1, $2, 1, $3, $4, $5)
		`
	}
	r, err := d.db.ExecContext(
		ctx,
		query,
		request.Name,
		request.Request.Body,
		time.Now().Unix(),
		request.Request.Method,
		request.CollectionId,
	)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		log.Println("Error fetching last row id: " + err.Error())
		return 0, err
	}

	err = insertHeadersValues(ctx, d.db, 0, id, request.Request.Headers, true)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *Database) GetRequests(ctx context.Context, collectionId int64) ([]models.RequestEntry, error) {
	query := `
		SELECT id, name, request_body, created_at
		FROM requests
		WHERE collection_id is null
	`

	if collectionId > 0 {
		query = `
			SELECT id, name, request_body, created_at
			FROM requests
			WHERE collection_id = $1
		`
	}
	r, err := d.db.QueryContext(
		ctx,
		query,
		collectionId,
	)
	if err != nil {
		log.Println("Error while getting requests" + err.Error())
		return nil, err
	}
	res := make([]models.RequestEntry, 0)
	for r.Next() {
		e := models.RequestEntry{}
		err = r.Scan(&e.Id, &e.Name, &e.Request.Body, &e.CreatedAt)
		if err != nil {
			log.Println("Error while scanning request fields" + err.Error())
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (d *Database) DeleteRequest(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, "DELETE FROM requests WHERE id = $1", id)
	if err != nil {
		log.Println("Error while executing the query...")
	}
	return err
}

func (d *Database) UpdateRequest(ctx context.Context, id int64, req *models.Request) (*models.StoredRequest, error) {
	_, err := d.db.ExecContext(ctx,
		`
		UPDATE requests
		SET name = $2,
		  	request_body = $3,
		  	method = $4
		WHERE id = $1
		`,
		id,
		req.Name,
		req.Request.Body,
		req.Request.Method,
	)
	if err != nil {
		return nil, err
	}
	if err = updateRequestHeadersValues(ctx, d.db, id, req.Request.Headers); err != nil {
		return nil, err
	}
	return &models.StoredRequest{Id: int(id), Name: req.Name, Request: req.Request}, nil
}

func (d *Database) GetRequest(ctx context.Context, id int64) (*models.StoredRequest, error) {
	r := d.db.QueryRowContext(
		ctx,
		`
		SELECT id, name, request_body, created_at
		FROM requests
		WHERE id = $1
		`,
		id,
	)
	req := models.StoredRequest{}
	if err := r.Scan(&req.Id, &req.Name, &req.Request.Body); err != nil {
		return nil, err
	}
	reqHeaders, _, err := getRequestHeaders(ctx, d.db, 0, id)
	if err != nil {
		return nil, err
	}
	req.Request.Headers = reqHeaders
	return &req, nil
}
