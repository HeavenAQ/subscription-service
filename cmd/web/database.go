package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Failed to connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Print("Error:", err, "Retrying to connect to database...")
		} else {
			log.Print("connected to database!")
			return connection
		}
		if counts > 10 {
			return nil
		}
		counts++

		log.Print("retrying to connect to database...")
		time.Sleep(1 * time.Second)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// check that db can be opened
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// check that db can be pinged
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}
