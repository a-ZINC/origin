package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"origin.me/internal/models"
	"origin.me/internal/store"
)

type SeriesHandler struct{ store *store.Store }

func NewSeriesHandler(st *store.Store) *SeriesHandler {
	return &SeriesHandler{store: st}
}

func (h *SeriesHandler) List(w http.ResponseWriter, r *http.Request) {
	series, err := h.store.ListSeries(r.Context())
	if err != nil {
		fail(w, "db error", 500)
		return
	}
	ok(w, series)
}

func (h *SeriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	sr, err := h.store.GetSeriesBySlug(r.Context(), slug)
	if err != nil || sr == nil {
		fail(w, "not found", 404)
		return
	}
	posts, _ := h.store.ListPostBySeries(r.Context(), sr.ID, IsAdminCtx(r.Context()))
	ok(w, map[string]interface{}{"series": sr, "posts": posts})
}

func (h *SeriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sr models.Series
	if err := readJson(r, &sr); err != nil {
		fail(w, "invalid body", 400)
		return
	}
	if sr.Slug == "" {
		sr.Slug = slugify(sr.Name)
	}
	if err := h.store.CreateSeries(r.Context(), &sr); err != nil {
		fail(w, err.Error(), 500)
		return
	}
	ok(w, sr)
}

func (h *SeriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	var sr models.Series
	if err := readJson(r, &sr); err != nil {
		fail(w, "invalid body", 400)
		return
	}
	sr.ID = chi.URLParam(r, "id")
	if err := h.store.UpdateSeries(r.Context(), &sr); err != nil {
		fail(w, err.Error(), 500)
		return
	}
	ok(w, sr)
}

func (h *SeriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.store.DeleteSeries(r.Context(), chi.URLParam(r, "id"))
	ok(w, map[string]bool{"deleted": true})
}