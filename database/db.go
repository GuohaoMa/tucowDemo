package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "host=db port=5432 user=postgres dbname=mydb password=pwd sslmode=disable")
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		panic(errors.Wrap(err, "error creating driver: %w"))
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver)
	if err != nil {
		panic(errors.Wrap(err, "error while migrating database: %w"))
	}

	if err = m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("no change made by migrations")
		} else {
			panic(errors.Wrap(err, "error while migrating database: %w"))
		}
	}
}
