package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products (
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func TestMain(m *testing.M) {
	fi, err := os.Open("dbinfo.txt")
	HandleErr(err)
	defer fi.Close()
	info := []string{}
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		info = append(info, scanner.Text())
	}
	a.Initalize(info[0], info[1], info[2], info[3], info[4])
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	_, err := a.Db.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.Db.Exec("DELETE FROM products")
	a.Db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("Get", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Expected an empty array. got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Fetch a Non-existent Product
func TestGetNotExistentProduct(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found. Got %s", m["error"])
	}
}

// Create a Product
func TestCreateProduct(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}
	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

// Fetch a Product
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct(1)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addProduct(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.Db.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

// // Update a Product
// func TestUpdateProduct(t *testing.T) {
// 	clearTable()
// 	addProduct(1)

// }
