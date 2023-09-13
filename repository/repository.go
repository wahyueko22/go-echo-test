// This file contains the repository implementation layer.
package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	//db, err := sql.Open("postgres", `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`)
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		panic(err)
	}
	return &Repository{
		Db: db,
	}
}
