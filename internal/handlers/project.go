package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"origin.me/internal/models"
	"origin.me/internal/store"
)

type ProjectHandler struct{ store *store.Store }

func NewProjectHandler(st *store.Store) *ProjectHandler {
	return &ProjectHandler{store: st}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.ListProjects(r.Context())
	if err != nil {
		fail(w, "db error", 500)
		return
	}
	ok(w, projects)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var pr models.Project
	if err := readJson(r, &pr); err != nil {
		fail(w, "invalid body", 400)
		return
	}
	if err := h.store.CreateProject(r.Context(), &pr); err != nil {
		fail(w, err.Error(), 500)
		return
	}
	ok(w, pr)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pr models.Project
	if err := readJson(r, &pr); err != nil {
		fail(w, "invalid body", 400)
		return
	}
	pr.ID = chi.URLParam(r, "id")
	if err := h.store.UpdateProject(r.Context(), &pr); err != nil {
		fail(w, err.Error(), 500)
		return
	}
	ok(w, pr)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.store.DeleteProject(r.Context(), chi.URLParam(r, "id"))
	ok(w, map[string]bool{"deleted": true})
}