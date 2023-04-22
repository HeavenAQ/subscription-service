package main

import (
	"log"
	"os"
	"sync"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initDB()
	db.Ping()

	// create sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime)

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &sync.WaitGroup{},
	}

	// set up mail

	// lisetn for web connections
}
