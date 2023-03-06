package infra

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "prueba_db"
)

var DbLocal *sql.DB

func ConnectDBLocal() *sql.DB {
	psqlconn := fmt.Sprintf("host= %s port= %d user= %s password= %s dbname= %s sslmode=disable", host, port, user, password, dbname)
	conection, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err.Error())
	}
	return conection
}

func CheckDBLocal() error {
	var err error
	ctx := context.Background()
	err = DbLocal.PingContext(ctx)
	if err != nil {
		return err
	}
	return err
}
