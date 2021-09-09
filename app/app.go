package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// connect database
func (a *App) Initalize(user, password, host, port, dbname string) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	a.Db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
}

// start this application
func (a *App) Run(addr string) {}
