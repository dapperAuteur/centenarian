package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	h := http.HandlerFunc(Echo)

	log.Println("listening on PORT: 8000")

	if err := http.ListenAndServe("localhost:8000", h); err != nil {
		log.Fatal(err)
	}
}

// Echo tells you about the request made
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You asked to", r.Method, r.URL.Path)
}
