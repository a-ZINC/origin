package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"origin.me/internal/handlers"
	"origin.me/internal/middleware"
	"origin.me/internal/store"
)

func renderTmpl(w http.ResponseWriter, tmpl *template.Template, name string, data map[string]interface{}) {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template: "+err.Error(), 500)
	}
}

func main() {
	dbUrl := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/testDb?sslmode=disable")
	jwtSecret := getenv("JWT_SECRET", "secret")
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

	authH := handlers.NewAuthHandler(st, jwtSecret)
	postH := handlers.NewPostHandler(st)
	projectH := handlers.NewProjectHandler(st)
	seriesH := handlers.NewSeriesHandler(st)
	analyticH := handlers.NewAnalyticsHandler(st)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.JWTMiddleware(jwtSecret))
	
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Post("/api/auth/login", authH.Login)
	r.Post("/api/track", analyticH.Track)

	r.Get("/api/posts", postH.List)
	r.Get("/api/posts/{slug}", postH.Get)
	r.Post("/api/posts/preview", postH.Preview)

	r.Get("/api/projects", projectH.List)

	r.Get("/api/series", seriesH.List)
	r.Get("/api/series/{slug}", seriesH.Get)


	r.Route("/api/admin", func(r chi.Router) {
		r.Use(middleware.RequireAdmin)

		r.Get("/posts", postH.ListAdmin)
		r.Post("/posts", postH.Create)
		r.Put("/posts/{id}", postH.Update)
		r.Delete("/posts/{id}", postH.Delete)

		r.Post("/series", seriesH.Create)
		r.Put("/series/{id}", seriesH.Update)
		r.Delete("/series/{id}", seriesH.Delete)

		r.Post("/projects",          projectH.Create)
		r.Put("/projects/{id}",      projectH.Update)
		r.Delete("/projects/{id}",   projectH.Delete)

		r.Get("/analytics",          analyticH.Overview)
		r.Get("/analytics/{id}",     analyticH.PostDetail)
	})

	log.Printf("API running → http://localhost:%s", port)
	srv := &http.Server{
		Addr: ":" + port,
		Handler: r,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func getenv(key ,de string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return de
}