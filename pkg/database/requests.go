package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/impulse-http/local-backend/pkg/models"
)

func insertHeadersValues(ctx context.Context, db *sql.DB, requestHistoryId, requestId int64, headers map[string][]string, isRequest bool) error {
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

func (d *Database) CreateRequest(ctx context.Context, request *models.NewRequestRequest) (int64, error) {
	r, err := d.db.ExecContext(ctx,
		`
		INSERT INTO requests(name, url, request_body, user_id, created_at, method) VALUES ($1, $2, 1, $3, $4)
	`, request.Name, request.Request.Url, request.Request.Body, time.Now().Unix(), request.Request.Method)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		log.Println("Error fetching last row id: " + err.Error())
		return 0, err
	}

	err = insertHeadersValues(ctx, d.db, 0, id, request.Request.Headers, false)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *Database) DeleteRequest(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx,
		`
			DELETE FROM requests WHERE id = $1
	`, id)
	if err != nil {
		return errors.Wrap(err, "delete request")
	}

	return nil
}

func (d *Database) GetListRequests(ctx context.Context) ([]*models.NewRequestRequest, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT id, name, request_body, user_id, method FROM requests`)
	if err != nil {
		return nil, errors.Wrap(err, "error while select")
	}
	ret := make([]*models.NewRequestRequest, 0)
	for rows.Next() {
		item := &models.NewRequestRequest{}
		err := rows.Scan(
			&item.Request.Id,
			&item.Name,
			&item.Request.Body,
			&item.Request.Id,
			&item.Request.Method,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}

	return ret, nil
}
