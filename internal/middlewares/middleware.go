package middlewares

import (
	"context"
	"go-chat/internal/utils"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == nil {
			claims, err := utils.ParseToken(cookie.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), "userID", claims.UserID)
				ctx = context.WithValue(ctx, "username", claims.Username)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
