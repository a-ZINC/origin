package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "home page")
	})

	r.Get("/blog/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		fmt.Fprintf(w, "you want post: %s", slug)
	})

	r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "form received")
	})

	fmt.Println("server started")

	http.ListenAndServe(":8080", r)
}
