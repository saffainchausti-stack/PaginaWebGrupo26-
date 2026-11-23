package handlers

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./static/index.html")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./static/about.html")
}
