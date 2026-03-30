package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"origin.me/internal/store"
)

type AuthHandler struct {
	store *store.Store
	jwtSceret []byte
}

func NewAuthHandler(store *store.Store, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		store: store,
		jwtSceret: []byte(jwtSecret),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJson(r, &body); err != nil {
		fail(w, "invalid json", 400)
		return
	}
	log.Println(body)

	user, err := h.store.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		fail(w, "db error", 500)
		return
	}
	log.Println(user)

	if user == nil {
		fail(w, "user not found", 404)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		fail(w, "invalid password", 401)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"isAdmin": true,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString(h.jwtSceret)
	if err != nil {
		fail(w, "failed to sign token", 500)
		return
	}
	ok(w, map[string]string{"token": signed, "email": user.Email})
}