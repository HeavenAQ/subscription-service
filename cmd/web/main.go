package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/HeavenAQ/subscription-service/data"
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
		Wait:     &wg,
		Models:   data.New(db),
	}

	// set up mail

	// gracefully shutdown the application
	go app.listenForShutdown()

	// listen for web connections
	app.serve()
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	app.InfoLog.Println("Cleaning up...")
	app.Wait.Wait()
	app.InfoLog.Println("Graceful shutdown complete")
}
