package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"origin.me/internal/handlers"
	"origin.me/internal/middleware"
	"origin.me/internal/models"
	render_ "origin.me/internal/render"
	"origin.me/internal/store"
)

func renderTmpl(w http.ResponseWriter, tmpl *template.Template, name string, data map[string]interface{}) {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template: "+err.Error(), 500)
	}
}

func main() {
	dbUrl := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/testDb?sslmode=disable")
	secret := getenv("SESSION_SECRET", "secret")
	adminEmail := getenv("ADMIN_EMAIL", "admin")
	adminPass := getenv("ADMIN_PASSWORD", "")
	port := getenv("PORT", "8080")

	st, err := store.New(dbUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("database connected")

	if adminPass != "" {
		if err := st.SeedAdmin(context.Background(), adminEmail, "ajinkya", adminPass); err != nil {
			log.Printf("seed admin: %v", err)
		}
	}

	authHandler := handlers.NewAuthHandler(st, sessStor)
	adminHandler := handlers.NewAdminHandler(st, sessStore)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RequestID)
	r.Use(func(next http.Handler) http.Handler {
		return middleware.SetAdmin(sessStore, next)
	})
	
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		posts, _ := st.ListPosts(r.Context(), middleware.IsAdmin(r.Context()))
		projects, _ := st.ListProjects(r.Context())

		fmt.Printf("post: %v, projects: %v")
	})

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")
		var posts []*models.Post
		if tag != "" {
			posts, _ = st.ListPostsByTag(r.Context(), tag, middleware.IsAdmin(r.Context()))
		} else {
			posts, _ = st.ListPosts(r.Context(), middleware.IsAdmin(r.Context()))
		}
	})

	r.Get("/blog/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		post, err := st.GetPostBySlug(r.Context(), slug)
		if err != nil || post == nil {
			http.NotFound(w, r)
			return
		}

		post.HTMLBody = render_.Markdown(post.Body)
	})

	r.Get("/series", func(w http.ResponseWriter, r *http.Request) {
		series, _ := st.ListSeries(r.Context())
	})

	r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		projects, _ := st.ListProjects(r.Context())
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
	})

	r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "form received")
	})

	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.LoginPost)
	r.Get("/logout", authHandler.Logout)

	r.Route("/admin", func(r chi.Router) {
		r.Use(func (next http.Handler) http.Handler  {
			return middleware.RequireAdmin(sessStore, next)
		})

		r.Get("/", adminHandler.Dashboard)
		r.Post("/preview", adminHandler.Preview)

		r.Get("/posts/new", adminHandler.NewPost)
		r.Post("/posts", adminHandler.CreatePost)
		r.Get("/posts/{id}/edit", adminHandler.EditPostPage)
		r.Post("/posts/{id}", adminHandler.UpdatePost)
		r.Post("/posts/{id}/delete", adminHandler.DeletePost)

		r.Get("/series/new", adminHandler.NewSeriesPage)
		r.Post("/series", adminHandler.CreateSeries)
		r.Get("/series/{id}/edit", adminHandler.EditSeriesPage)
		r.Post("/series/{id}", adminHandler.UpdateSeries)
		r.Post("/series/{id}/delete", adminHandler.DeleteSeries)

		r.Get("/projects/new", adminHandler.NewProjectPage)
		r.Post("/projects", adminHandler.CreateProject)
		r.Get("/projects/{id}/edit", adminHandler.EditProjectPage)
		r.Post("/projects/{id}", adminHandler.UpdateProject)
		r.Post("/projects/{id}/delete", adminHandler.DeleteProject)
	})

	fmt.Println("server started")

	http.ListenAndServe(":8080", r)
}

func getenv(key ,de string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return de
}