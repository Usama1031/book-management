package middleware

import (
	"context"
	"net/http"

	"github.com/usama1031/book-management/pkg/helpers"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authorization Required!", http.StatusUnauthorized)
			return
		}

		clientToken := cookie.Value

		claims, msg := helpers.ValidateToken(clientToken)

		if msg != "" {
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "first_name", claims.First_name)
		ctx = context.WithValue(ctx, "last_name", claims.Last_name)
		ctx = context.WithValue(ctx, "uid", claims.Uid)
		ctx = context.WithValue(ctx, "user_type", claims.User_type)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
