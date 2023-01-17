package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var counst int64

const port = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting auth service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to DB")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet...")
			counst++
		} else {
			log.Println("Postgres ready")
			return connection
		}

		if counst > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Bacing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
