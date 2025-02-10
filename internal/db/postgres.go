package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DB interface {
	MustBegin() Tx
}

type Postgres struct {
	*sqlx.DB
}

func NewPostgres(host, username, password, dbName string) *Postgres {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbName)
	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		log.New().Fatal(err)
	}
	return &Postgres{db}
}

func (p *Postgres) MustBegin() Tx {
	return p.DB.MustBegin()
}
