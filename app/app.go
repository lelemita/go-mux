package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// connect database
func (a *App) Initialize(user, password, host, port, dbname string) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	fmt.Println(connectionString)
	a.Db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

func respWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp, _ := json.Marshal(payload)
	w.Write(resp)
}

func respWithError(w http.ResponseWriter, code int, message string) {
	payload := map[string]string{"error": message}
	respWithJSON(w, code, payload)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	prod := product{ID: id}
	err = prod.getProduct(a.Db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respWithError(w, http.StatusNotFound, "Product not found")
		default:
			respWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respWithJSON(w, http.StatusOK, prod)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	products, err := getProducts(a.Db, offset, limit)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err = p.createProduct(a.Db)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&p)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	err = p.updateProduct(a.Db)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	err = p.deleteProduct(a.Db)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// start this application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}
