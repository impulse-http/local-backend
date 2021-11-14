package database

import (
	"database/sql"
	"fmt"
	"github.com/impulse-http/local-backend/pkg"
	"log"
	"net/http"
	"time"
)

func insertHeadersValues(db *sql.DB, rid int64, headers map[string][]string, isRequest bool) error {
	for name, values := range headers {
		r, err := db.Exec(
			`
			INSERT INTO headers (key, request_id, is_request)
			VALUES ($1, $2, $3);
			`,
			name,
			rid,
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
			r, err = db.Exec(
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

func (d *Database) CreateHistoryEntry(req *pkg.RequestType, res *pkg.ResponseType) (int64, error) {
	db := d.db
	r, err := db.Exec(
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

	err = insertHeadersValues(db, id, req.Headers, true)
	if err != nil {
		return id, err
	}

	err = insertHeadersValues(db, id, res.Headers, false)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (d *Database) GetHistory() ([]RequestHistoryEntry, error) {
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

	newHistoryEntry := func() RequestHistoryEntry {
		entry := RequestHistoryEntry{}
		entry.Request.Headers = make(http.Header)
		entry.Response.Headers = make(http.Header)
		return entry
	}

	db := d.db
	rows, err := db.Query(`
			SELECT rh.id, rh.method, rh.request_body, rh.response_body, rh.created_at,
			       h.is_request, h.key, hv.header_value
			FROM requests_history rh
			    JOIN headers h on rh.id = h.request_id
				JOIN headers_values hv ON hv.header_id = h.id
			ORDER BY rh.id, h.id;
	`)

	if err != nil {
		log.Println("Request history query failed" + fmt.Sprint(err))
		return nil, err
	}

	res := make([]RequestHistoryEntry, 0)
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
