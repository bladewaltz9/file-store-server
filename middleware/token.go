package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/bladewaltz9/file-store-server/models"
	"github.com/dgrijalva/jwt-go"
)

// extractToken: extract the token from the request
func extractToken(r *http.Request) string {
	// get the token from the header
	tokenStr := r.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// get the token from the cookie if the header is missing
	if tokenStr == "" {
		tokenCookie, err := r.Cookie("token")
		if err == nil {
			tokenStr = tokenCookie.Value
		}
	}
	return tokenStr
}

// validateToken: validate the token
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	// parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	// get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

// TokenAuthMiddleware: middleware to authenticate the token
func TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// set the claims to the context
		ctx := context.WithValue(r.Context(), models.ContextKey("claims"), claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// IsAuthenticated: check if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	tokenStr := extractToken(r)
	if tokenStr == "" {
		return false
	}

	_, err := ValidateToken(tokenStr)
	return err == nil
}
