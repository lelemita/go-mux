package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lelemita/go-mux/app"
)

const port string = ":8010"

var a = app.App{}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// the entry point for this application
func main() {
	fi, err := os.Open("dbinfo.txt")
	HandleErr(err)
	defer fi.Close()
	info := []string{}
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		info = append(info, scanner.Text())
	}
	a.Initialize(info[0], info[1], info[2], info[3], info[4])
	fmt.Printf("Listening on http://localhost%s\n", port)
	a.Run(port)
}
