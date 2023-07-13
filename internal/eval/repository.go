package eval

import (
	"context"
	"github.com/svilenkomitov/eval/internal/storage"
)

type Endpoint string

const (
	Validate Endpoint = "/validate"
	Evaluate Endpoint = "/evaluate"
)

type Error struct {
	Expression string    `json:"expression" db:"expression"`
	Endpoint   Endpoint  `json:"endpoint" db:"endpoint"`
	Frequency  int       `json:"frequency" db:"frequency"`
	Kind       ErrorKind `json:"kind" db:"kind"`
}

type Repository interface {
	Upsert(error Error) error
	Errors() ([]Error, error)
}

type repository struct {
	db *storage.Database
}

func NewRepository(db *storage.Database) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Upsert(error Error) error {
	query := `INSERT INTO errors (expression, endpoint, frequency, kind) 
				VALUES ($1, $2, 1, $3)
				ON CONFLICT ON constraint errors_unique_constraint
				    DO UPDATE SET frequency = errors.frequency + 1;`

	_, err := r.db.ExecContext(context.Background(), query,
		error.Expression, error.Endpoint, error.Kind)
	return err
}

func (r repository) Errors() ([]Error, error) {
	query := `SELECT expression, endpoint, frequency, kind FROM errors;`
	errors := []Error{}
	rows, err := r.db.Query(query)
	if err != nil {
		return []Error{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var error Error
		err := rows.Scan(&error.Expression, &error.Endpoint, &error.Frequency, &error.Kind)
		if err != nil {
			return []Error{}, err
		}
		errors = append(errors, error)
	}

	if err := rows.Err(); err != nil {
		return []Error{}, err
	}

	return errors, nil
}
