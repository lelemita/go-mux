package main

import (
	"bufio"
	"log"
	"os"
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
