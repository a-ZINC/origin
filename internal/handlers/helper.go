package handlers

import (
	"context"
	"strings"
)

type ctxKey string
const adminKey ctxKey = "isAdmin"

func IsAdminCtx(ctx context.Context) bool {
	val, ok := ctx.Value(adminKey).(bool)
	if !ok {
		return false
	}
	return val
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

func WithAdmin(ctx context.Context, val bool) context.Context {
	return context.WithValue(ctx, adminKey, val)
}

