/*
Package api exposes the main API engine. All HTTP APIs are handled here - so-called "business logic" should be here, or
in a dedicated package (if that logic is complex enough).

To use this package, you should create a new instance with New() passing a valid Config. The resulting Router will have
the Router.Handler() function that returns a handler that can be used in a http.Server (or in other middlewares).

Example:

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: appdb,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	// ... other stuff here, like middleware chaining, etc.

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()

See the `main.go` file inside the `cmd/webapi` for a full usage example.
*/
package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Shrek-the-ogre19/WASAText/service/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database.AppDatabase where data are saved
	Database database.AppDatabase
}

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Close terminates any resource used in the package
	Close() error
}

// New returns a new Router instance
func New(cfg Config) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}

type _router struct {
	router *httprouter.Router

	// baseLogger is a logger for non-requests contexts, like goroutines or background tasks not started by a request.
	// Use context logger if available (e.g., in requests) instead of this logger.
	baseLogger logrus.FieldLogger

	db database.AppDatabase
}

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
