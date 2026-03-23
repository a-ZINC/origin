package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"origin.me/internal/store"
)

func main() {

	st, err := store.New("postgres://postgres:postgres@localhost:5432/testDb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	fmt.Println("database connected")

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "home page")
	})

	r.Get("/blog/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		
		post, err := st.GetPostBySlug(r.Context(), slug)
		if err != nil {
			fmt.Printf("err: %v", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		if post == nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "post: %s", post.Title)
	})

	r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "form received")
	})

	fmt.Println("server started")

	http.ListenAndServe(":8080", r)
}
