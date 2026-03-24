package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"origin.me/internal/models"
	"origin.me/internal/store"
	render_ "origin.me/internal/render"
)

func render(w http.ResponseWriter, tmpl *template.Template, name string, data map[string]interface{}) {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	st, err := store.New("postgres://postgres:postgres@localhost:5432/testDb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	fmt.Println("database connected")

	tmpl, err := loadTemplates()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")
		var posts []*models.Post
		var err error

		if tag != "" {
			posts, err = st.ListPostsByTag(r.Context(), tag, false)
		} else {
			posts, err = st.ListPosts(r.Context(), false)
		}

		if err != nil {
			fmt.Printf("err: %v", err)
			http.Error(w, "something broke", http.StatusInternalServerError)
			return
		} 

		render(w, tmpl, "home.html", map[string]interface{}{
			"Posts":   posts,
			"IsAdmin": false,
		})
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

		post.HTMLBody = render_.Markdown(post.Body)

		render(w, tmpl, "post.html", map[string]interface{}{
			"Post":    post,
			"IsAdmin": false,
		})
	})

	r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "form received")
	})

	fmt.Println("server started")

	http.ListenAndServe(":8080", r)
}

func loadTemplates() (*template.Template, error) {
	funcMap := template.FuncMap{
		"formatDate": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("Jan 2, 2006")
		},
		"safeHTML": func(s interface{}) template.HTML {
			if s == nil {
				return ""
			}
			switch v := s.(type) {
			case string:
				return template.HTML(v)
			case []byte:
				return template.HTML(string(v))
			default:
				return ""
			}
		},
		"join": strings.Join,
	}

	tmpl := template.New("").Funcs(funcMap)

	err := filepath.Walk("templates", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		_, err = tmpl.ParseFiles(path)
		return err
	})

	return tmpl, err
}
