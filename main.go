package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lelemita/go-mux/app"
)

const port = ":8010"

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// the entry point for this application
func main() {
	a := app.App{}
	bytes, err := ioutil.ReadFile("dbinfo.txt")
	HandleErr(err)
	strs := fmt.Sprint(string(bytes))
	info := strings.Split(strs, "\n")
	a.Initalize(os.Getenv(info[0]), os.Getenv(info[0]), os.Getenv(info[2]))
	fmt.Printf("Listening on http://localhost%s\n", port)
	a.Run(port)
}
