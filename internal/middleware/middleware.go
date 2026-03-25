package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

type ctxKey string
const adminKey ctxKey = "isAdmin"

func SetAdmin(store sessions.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := store.Get(r, "session")
		isAdmin, _ := (sess.Values["admin"]).(bool)

		ctx := context.WithValue(r.Context(), adminKey, isAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAdmin(store sessions.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := r.Context().Value(adminKey).(bool)
		if !isAdmin {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IsAdmin(ctx context.Context) bool {
	return ctx.Value(adminKey).(bool)
}