package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/furkankorkmaz309/todo-api/internal/app"
	"github.com/furkankorkmaz309/todo-api/internal/db"
	"github.com/furkankorkmaz309/todo-api/internal/routes"
)

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)

	db, err := db.InitDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &app.App{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		DB:       db,
	}

	addr := flag.String("addr", ":8080", "new http port")
	flag.Parse()

	router := routes.Routes(app)

	srv := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	app.InfoLog.Println("Server running on port", *addr)
	go func() {
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.ErrorLog.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	select {} // ?
}
