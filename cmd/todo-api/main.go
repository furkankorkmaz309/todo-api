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

	app := &app.App{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		DB:       db,
	}

	if err != nil {
		app.ErrorLog.Fatal(err)
	}
	defer db.Close()

	addr := flag.String("addr", ":8080", "new http port")

	mux := routes.Routes(app)

	app.InfoLog.Println("Server running on port", *addr)
	http.ListenAndServe(*addr, mux)
}
