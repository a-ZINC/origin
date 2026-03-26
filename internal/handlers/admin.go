package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"origin.me/internal/models"
	"origin.me/internal/render"
	"origin.me/internal/store"
)

type AdminHandler struct {
	store *store.Store
	tmpl *template.Template
	sessions sessions.Store
}

func NewAdminHandler(st *store.Store, tmpl *template.Template, ss sessions.Store) *AdminHandler {
	return &AdminHandler{store: st, tmpl: tmpl, sessions: ss}
}

func (h *AdminHandler) render(w http.ResponseWriter, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["IsAdmin"] = true
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template: "+err.Error(), 500)
	}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.ListAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	series, _ := h.store.ListSeries(r.Context())
	projects, _ := h.store.ListProjects(r.Context())

	h.render(w, "admin_dashboard.html", map[string]interface{}{
		"Posts": posts,
		"Series": series,
		"Projects": projects,
	})
}

// ---- Post ----

func (h *AdminHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	series, _ := h.store.ListSeries(r.Context())

	h.render(w, "admin_new_post.html", map[string]interface{}{
		"Post": &models.Post{Visibility: "public"},
		"Series": series,
		"IsNew": true,
	})
}

func (h *AdminHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	p := postFromForm(r)
	if err := h.store.CreatePost(r.Context(), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AdminHandler) EditPostPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.store.GetPostByID(r.Context(), id)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}
	series, _ := h.store.ListSeries(r.Context())
	h.render(w, "admin_post_form.html", map[string]interface{}{
		"Post":   post,
		"Series": series,
		"IsNew":  false,
	})
}

func (h *AdminHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p := postFromForm(r)
	p.ID = id

	existing, _ := h.store.GetPostByID(r.Context(), id)
	if existing != nil {
		p.UpdatedAt = existing.UpdatedAt
	}

	if p.UpdatedAt.IsZero()  {
		p.UpdatedAt = time.Now()
	}

	if err := h.store.UpdatePost(r.Context(), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AdminHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.store.DeletePost(r.Context(), id)
	http.Redirect(w, r, "/admin", http.StatusFound)
}


// --- Series ---

func (h *AdminHandler) NewSeriesPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, "admin_series_form.html", map[string]interface{}{
		"Series": &models.Series{},
		"IsNew":  true,
	})
}

func (h *AdminHandler) CreateSeries(w http.ResponseWriter, r *http.Request) {
	sr := &models.Series{
		Name:        r.FormValue("name"),
		Slug:        slugify(r.FormValue("slug")),
		Description: r.FormValue("description"),
	}
	if sr.Slug == "" {
		sr.Slug = slugify(sr.Name)
	}
	if err := h.store.CreateSeries(r.Context(), sr); err != nil {
		http.Error(w, "create series: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AdminHandler) EditSeriesPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sr, err := h.store.GetSeriesByID(r.Context(), id)
	if err != nil || sr == nil {
		http.NotFound(w, r)
		return
	}
	h.render(w, "admin_series_form.html", map[string]interface{}{
		"Series": sr,
		"IsNew":  false,
	})
}

func (h *AdminHandler) UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sr := &models.Series{
		ID:          id,
		Name:        r.FormValue("name"),
		Slug:        slugify(r.FormValue("slug")),
		Description: r.FormValue("description"),
	}
	if err := h.store.UpdateSeries(r.Context(), sr); err != nil {
		http.Error(w, "update series: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AdminHandler) DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.store.DeleteSeries(r.Context(), id)
	http.Redirect(w, r, "/admin", http.StatusFound)
}


// --- Projects ---

func (h *AdminHandler) NewProjectPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, "admin_project_form.html", map[string]interface{}{
		"Project": &models.Project{},
		"IsNew":   true,
	})
}

func (h *AdminHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	pr := projectFromForm(r)
	if err := h.store.CreateProject(r.Context(), pr); err != nil {
		http.Error(w, "create project: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AdminHandler) EditProjectPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	projects, _ := h.store.ListProjects(r.Context())
	var found *models.Project
	for _, p := range projects {
		if p.ID == id {
			found = p
			break
		}
	}
	if found == nil {
		http.NotFound(w, r)
		return
	}
	h.render(w, "admin_project_form.html", map[string]interface{}{
		"Project": found,
		"IsNew":   false,
	})
}

func (h *AdminHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pr := projectFromForm(r)
	pr.ID = id
	if err := h.store.UpdateProject(r.Context(), pr); err != nil {
		http.Error(w, "update project: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AdminHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.store.DeleteProject(r.Context(), id)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// --- Helpers ---

func postFromForm(r *http.Request) *models.Post {
	r.ParseForm()
	p := &models.Post{
		Title: r.FormValue("title"),
		Slug:       slugify(r.FormValue("slug")),
		Excerpt:    r.FormValue("excerpt"),
		Body:       r.FormValue("body"),
		CoverImage: r.FormValue("cover_image"),
		Visibility: r.FormValue("visibility"),
		Published:  r.FormValue("published") == "true",
	}

	if p.Visibility == "" {
		p.Visibility = "public"
	}
	if p.Slug == "" {
		p.Slug = slugify(p.Title)
	}
	
	for _, t := range strings.Split(r.FormValue("tags"), ",") {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" {
			p.Tags = append(p.Tags, t)
		}
	}

	if sid := r.FormValue("seriesId"); sid != "" {
		p.SeriesId = &sid
	}
	if sp := r.FormValue("seriesPos"); sp != "" {
		if n, err := strconv.Atoi(sp); err == nil {
			p.SeriesPos = &n
		}
	}
	return p
}

func projectFromForm(r *http.Request) *models.Project {
	r.ParseForm()
	stars, _ := strconv.Atoi(r.FormValue("stars"))
	order, _ := strconv.Atoi(r.FormValue("sortOrder"))

	return &models.Project{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		URL:         r.FormValue("url"),
		RepoURL:     r.FormValue("repoUrl"),
		Language:    r.FormValue("language"),
		Stars:       stars,
		Featured:    r.FormValue("featured") == "true",
		SortOrder:   order,
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder

	for _, ch := range s {
		switch {
			case ch >= '0' && ch <= '9', ch >= 'a' && ch <= 'z':
				b.WriteRune(ch)
			case ch == ' ' || ch == '-' || ch == '_':
				b.WriteRune('-')
		}
	}
	slug := strings.Trim(b.String(), "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return slug
}

func (h *AdminHandler) Preview(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(render.Markdown(body)))
}

