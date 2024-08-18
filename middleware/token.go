package middleware

import (
	"context"
	"net/http"

	"github.com/bladewaltz9/file-store-server/models"
	"github.com/dgrijalva/jwt-go"
)

// TokenAuthMiddleware: middleware to authenticate the token
func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the header
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("file-store-server"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// get the claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claim", http.StatusUnauthorized)
			return
		}

		// set the claims to the context
		ctx := context.WithValue(r.Context(), models.ContextKey("claims"), claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
