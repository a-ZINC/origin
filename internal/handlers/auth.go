package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"origin.me/internal/store"
)

type AuthHandler struct {
	store *store.Store
	tmpl *template.Template
	sessions sessions.Store
}

func NewAuthHandler(store *store.Store, tmpl *template.Template, sessions sessions.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
		tmpl: tmpl,
		sessions: sessions,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
		"Error": r.URL.Query().Get("error"),
	})
}

func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.store.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Redirect(w, r, "/login?error=invalid+credentials", http.StatusFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=invalid+credentials", http.StatusFound)
		return
	}

	sess, err := h.sessions.Get(r, "session")
	sess.Values["admin"] = true
	sess.Save(r, w)

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.sessions.Get(r, "session")
	sess.Values["admin"] = false
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}