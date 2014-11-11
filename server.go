package main

import (
	"log"
	"net/http"
	"os"
)

func verbose(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func main() {
	var port string
	if len(os.Args) != 2 {
		port = "8000"
	} else {
		port = os.Args[1]
	}

	http.Handle("/", verbose(http.FileServer(http.Dir("www"))))

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}
