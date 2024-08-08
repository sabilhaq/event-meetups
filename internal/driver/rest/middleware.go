package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("event-maker")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	contextKeyUser = contextKey("user")
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Render(w, r, NewErrorResp(NewUnauthorizedError()))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			render.Render(w, r, NewErrorResp(NewUnauthorizedError()))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func UserFromContext(ctx context.Context) int {
	return ctx.Value(contextKeyUser).(int)
}
