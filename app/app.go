package app

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// connect database
func (a *App) Initalize(user, passwrd, dbname string) {}

// start this application
func (a *App) Run(addr string) {}
