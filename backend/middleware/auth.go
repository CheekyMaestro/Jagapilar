package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("JAGAPILAR_SECRET_KEY_CHANGE_ME_IN_PROD")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const UserIDKey contextKey = "userID"
const RoleKey contextKey = "role"

// Auth is a middleware that extracts JWT from Authorization header
// and puts UserID and Role into the request context.
// Endpoints can still process the request if token is missing,
// it's up to the handler to enforce if they need authentication.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err == nil && token.Valid {
				// Token is valid, put claims into context
				ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
				ctx = context.WithValue(ctx, RoleKey, claims.Role)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
