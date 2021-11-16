package database

import (
	"context"
	"github.com/impulse-http/local-backend/pkg/models"
	"log"
)

func (d *Database) CreateCollection(ctx context.Context, c models.Collection) (*models.StoredCollection, error) {
	r, err := d.db.ExecContext(
		ctx,
		`
		INSERT INTO collections (name)
		VALUES ($1)
		`,
		c.Name,
	)
	if err != nil {
		log.Println("Error while inserting to database...")
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		log.Println("Error while fetching last row id...")
		return nil, err
	}
	return &models.StoredCollection{Id: id, Name: c.Name}, nil
}

func (d *Database) DeleteCollection(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, "DELETE FROM collections WHERE id = $1", id)
	if err != nil {
		log.Println("Error while executing the query...")
	}
	return err
}

func (d *Database) ListCollections(ctx context.Context) ([]models.StoredCollection, error) {
	r, err := d.db.QueryContext(ctx, "SELECT id, name FROM collections;")
	if err != nil {
		log.Println("Error while getting the list of collections...")
		return nil, err
	}
	res := make([]models.StoredCollection, 0)
	for r.Next() {
		c := models.StoredCollection{}
		err = r.Scan(&c.Id, &c.Name)
		if err != nil {
			log.Println("Error while scanning collection...")
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (d *Database) UpdateCollection(ctx context.Context, id int64, c models.StoredCollection) (*models.StoredCollection, error) {
	_, err := d.db.ExecContext(ctx, "UPDATE collections SET name = $1 WHERE id = $2", c.Name, id)
	if err != nil {
		log.Println("Error while updating a collecton")
		return nil, err
	}
	return &models.StoredCollection{Id: id, Name: c.Name}, nil
}
