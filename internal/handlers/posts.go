package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"origin.me/internal/models"
	"origin.me/internal/render"
	"origin.me/internal/store"
)

type PostHandler struct {
	store *store.Store
}

func NewPostHandler(store *store.Store) *PostHandler {
	return &PostHandler{store: store}
}

func(h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	isAdmin := isAdminCtx(r.Context())
	tag := r.URL.Query().Get("tag")

	var posts interface{}
	var err error
	
	if tag != "" {
			posts, _ = h.store.ListPostsByTag(r.Context(), tag, isAdmin)
	} else {
		posts, _ = h.store.ListPosts(r.Context(), isAdmin)
	}

	if err != nil {
		fail(w, "failed to fetch posts", 500)
		return
	}
	ok(w, posts)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	post, err := h.store.GetPostBySlug(r.Context(), slug)
	if err != nil {
		fail(w, "db error", 500)
		return
	}
	if post == nil {
		fail(w, "not found", 404)
		return
	}

	isAdmin := IsAdminCtx(r.Context())
	if !post.Published && !isAdmin {
		fail(w, "not found", 404)
		return
	}
	if post.Visibility == "private" && !isAdmin {
		fail(w, "not found", 404)
		return
	}

	post.HTMLBody = render.Markdown(post.Body)
	ok(w, post)
}

func (h *PostHandler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.ListAllPosts(r.Context())
	if err != nil {
		fail(w, "db error", 500)
		return
	}
	ok(w, posts)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Body  string `json:"body"`
		Excerpt string `json:"excerpt"`
		CoverImage string `json:"coverImage"`
		Tags []string `json:"tags"`
		Published bool `json:"published"`
		Visibility string `json:"visibility"`
		SeriesId *string `json:"seriesId"`
		SeriesPos *int `json:"seriesPos"`
	}

	if err := readJson(r, &body); err != nil {
		fail(w, "invalid json", 400)
		return
	}

	p := &models.Post{
		Title: body.Title,
		Slug: body.Slug,
		Body: body.Body,
		Excerpt: body.Excerpt,
		CoverImage: body.CoverImage,
		Tags: body.Tags,
		Published: body.Published,
		Visibility: body.Visibility,
		SeriesId: body.SeriesId,
		SeriesPos: body.SeriesPos,
	}
	if p.Visibility == "" {
		p.Visibility = "public"
	}
	if p.Slug == "" {
		p.Slug = slugify(p.Title)
	}

	if err := h.store.CreatePost(r.Context(), p); err != nil {
		fail(w, "create failed: "+err.Error(), 500)
		return
	}
	ok(w, p)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Body  string `json:"body"`
		Excerpt string `json:"excerpt"`
		CoverImage string `json:"coverImage"`
		Tags []string `json:"tags"`
		Published bool `json:"published"`
		Visibility string `json:"visibility"`
		SeriesId *string `json:"seriesId"`
		SeriesPos *int `json:"seriesPos"`
	}

	if err := readJson(r, &body); err != nil {
		fail(w, "invalid json", 400)
		return
	}

	p := &models.Post{
		ID: id,
		Title: body.Title,
		Slug: body.Slug,
		Body: body.Body,
		Excerpt: body.Excerpt,
		CoverImage: body.CoverImage,
		Tags: body.Tags,
		Published: body.Published,
		Visibility: body.Visibility,
		SeriesId: body.SeriesId,
		SeriesPos: body.SeriesPos,
	}
	
	if err := h.store.UpdatePost(r.Context(), p); err != nil { 
		fail(w, "update failed: "+err.Error(), 500)
		return 
	}
	ok(w, p)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.store.DeletePost(r.Context(), id); err != nil {
		fail(w, "delete failed: "+err.Error(), 500)
		return
	}
	ok(w, nil)
}

func (h *PostHandler) Preview(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Markdown string `json:"markdown"`
	}
	if err := readJson(r, &body); err != nil {
		fail(w, "invalid body", 400)
		return
	}
	ok(w, map[string]string{"html": render.Markdown(body.Markdown)})
}

