package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

const jwtSigningKey = "WASAText-signing-key"

type authClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"uid"`
}

func (rt *_router) createToken(userID int) (string, error) {
	now := time.Now()
	claims := authClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSigningKey))
}

func (rt *_router) userIDFromToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(jwtSigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*authClaims)
	if !ok || !token.Valid {
		return 0, jwt.ErrTokenInvalidClaims
	}
	return claims.UserID, nil
}

func bearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}

func (rt *_router) authorize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (int, bool) {
	token, ok := bearerToken(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, false
	}

	userID, err := rt.userIDFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, false
	}

	pathUserID, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 0, false
	}

	if userID != pathUserID {
		w.WriteHeader(http.StatusForbidden)
		return 0, false
	}

	return userID, true
}
