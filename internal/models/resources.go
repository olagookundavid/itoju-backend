package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/olagookundavid/itoju/internal/validator"
)

type Resources struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	ImageUrl string   `json:"image_url"`
	Link     string   `json:"link"`
	Tags     []string `json:"tags"`
	Version  int32    `json:"-"`
}

type ResourcesModel struct {
	DB *sql.DB
}

func ValidateResource(v *validator.Validator, resource *Resources) {
	v.Check(resource.Name != "", "name", "must be provided")
	v.Check(resource.ImageUrl != "", "image url", "must be provided")
	v.Check(resource.Link != "", "link", "must be provided")

}

func (m ResourcesModel) GetResources() ([]*Resources, error) {
	query := ` SELECT * FROM resources `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []*Resources{}, err
	}
	defer rows.Close()
	resources := []*Resources{}
	for rows.Next() {
		var resource Resources
		err := rows.Scan(&resource.Id, &resource.Name, &resource.ImageUrl, &resource.Link, pq.Array(&resource.Tags), &resource.Version)
		if err != nil {
			return []*Resources{}, err
		}

		resources = append(resources, &resource)
	}
	if err = rows.Err(); err != nil {
		return []*Resources{}, err
	}
	return resources, nil
}

func (m ResourcesModel) InsertResources(resources Resources) error {
	query := `
	INSERT INTO resources (name, imageUrl, link, tags)
	VALUES ($1, $2, $3, $4) `

	args := []any{resources.Name, resources.ImageUrl, resources.Link, pq.Array(resources.Tags)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "resources_name_key"`:

			return ErrRecordAlreadyExist
		default:
			return err
		}
	}
	return err
}

func (m ResourcesModel) UpdateResources(resources *Resources) error {

	query := ` UPDATE resources SET name = $1, imageUrl = $2, link = $3, tags = $4, version = version + 1 WHERE id = $5 AND version = $6; `

	args := []any{resources.Name, resources.ImageUrl, resources.Link, pq.Array(resources.Tags), resources.Id, resources.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m ResourcesModel) Get(id int64) (*Resources, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := ` SELECT * FROM resources WHERE id = $1; `
	var resource Resources
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&resource.Id,
		&resource.Name,
		&resource.ImageUrl,
		&resource.Link,
		pq.Array(&resource.Tags),
		&resource.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &resource, nil
}

func (m ResourcesModel) DeleteResources(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := ` DELETE FROM resources WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
