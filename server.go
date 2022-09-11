package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
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

/* Scan out description of an image and return it */
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
	db, err := sql.Open("sqlite3", "./images.db")
	if err != nil {
		panic(err.Error())
	}

	/* Return data in database as json */
	http.HandleFunc("/sql", func(w http.ResponseWriter, r *http.Request) {
		/* Run query */
		rows, err := db.Query("SELECT * FROM images")
		if err != nil {
			panic(err.Error())
		}

		/* Provide correct header */
		w.Header().Set("Content-Type", "application/json")

		/* Extract rows and write to HTTP response */
		var img Image

		for rows.Next() {
			rows.Scan(&img.ID, &img.Name, &img.Desc, &img.URL)
			obj, err := json.Marshal(img)
			if err != nil {
				panic(err.Error())
			}
			w.Write(obj)
		}
	})

	/* Close database */
	defer db.Close()

	/* Run http server */
	log.Fatal(http.ListenAndServe(":8081", nil))

}
