package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"
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

func main() {

	/* General handler */
	http.Handle("/", http.FileServer(http.Dir("./html")))

	/* Handler for /hi */
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	/* Handler for incrementor */
	http.HandleFunc("/increment", incrementCounter)

	/* Run http server */
	log.Fatal(http.ListenAndServe(":8081", nil))

}
