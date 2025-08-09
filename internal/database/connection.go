package database

import (
	"goserv/ent/gen"

	_ "github.com/lib/pq"
)

func NewDB(dsn string) (*gen.Client, error) {
	client, err := gen.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return client, nil
}
