package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

/* Globals */
var counter int
var mutex = &sync.Mutex{}

/* echo path */
func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %q", html.EscapeString(r.URL.RawPath))
}

/* Increment counter and display on page */
func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
}

/* Structs for SQL extraction */
type Image struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
	Desc string `json:"description"`
	URL  string `json:"URL"`
}

func scanDescription(row *sql.Row) string {
	var img Image

	err := row.Scan(&img.ID, &img.Name, &img.Desc, &img.URL)
	if err != nil {
		panic(err.Error())
	}

	return img.Desc
}

func main() {

	/* General handler */
	http.Handle("/", http.FileServer(http.Dir("./html")))

	/* Handler for /hi */
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	/* Handler for incrementor */
	http.HandleFunc("/increment", incrementCounter)

	/* Open a database connection */
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Image_gallery")

	/* Print error */
	if err != nil {
		panic(err.Error())
	}

	/* Extract first image and print description */
	http.HandleFunc("/sql", func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRow("SELECT ID, name, description, URL FROM image WHERE ID=?", 1)

		fmt.Fprintf(w, scanDescription(row))
	})

	/* Close database */
	defer db.Close()

	/* Run http server */
	log.Fatal(http.ListenAndServe(":8081", nil))

}
