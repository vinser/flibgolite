package opds

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os/user"
	"strings"
	"time"

	"github.com/kataras/basicauth"
	"golang.org/x/crypto/bcrypt"
)

type Middleware func(http.Handler) http.Handler

// func (h *Handler) NewBasicAuth() basicauth.Middleware {
func (h *Handler) NewAuth() Middleware {
	switch h.CFG.Auth.METHOD {
	case "plain":
		return h.NewBasicAuthPlain()
	case "file":
		return h.NewBasicAuthFile()
	case "db":
		return h.NewBasicAuthDB()
	default:
		return h.NewNoneAuth()
	}
}

func newBasicAuth(allowFunc func(r *http.Request, username, password string) (interface{}, bool)) basicauth.Middleware {
	opts := basicauth.Options{
		Realm:        basicauth.DefaultRealm,
		ErrorHandler: basicauth.DefaultErrorHandler,
		MaxAge:       2 * time.Hour,
		GC: basicauth.GC{
			Every: 3 * time.Hour,
		},
		Allow: allowFunc,
	}
	return basicauth.New(opts)
}

func allowUser(username, password, expectedUsername, expectedPassword string, useBcrypt bool) bool {
	var usernameMatch, passwordMatch bool
	usernameHash := sha256.Sum256([]byte(username))
	expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
	usernameMatch = (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
	if !useBcrypt {
		expectedPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(expectedPassword), bcrypt.MinCost)
		passwordMatch = bcrypt.CompareHashAndPassword(expectedPasswordHash, []byte(password)) == nil
	} else {
		passwordHash := sha256.Sum256([]byte(password))
		expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))
		passwordMatch = (subtle.ConstantTimeCompare(expectedPasswordHash[:], passwordHash[:]) == 1)
	}
	return usernameMatch && passwordMatch
}
func (h *Handler) NewNoneAuth() Middleware {
	return func(next http.Handler) http.Handler {
		return next
	}
}

func (h *Handler) NewBasicAuthPlain() basicauth.Middleware {
	return newBasicAuth(h.AllowUserPlain)
}

func (h *Handler) AllowUserPlain(r *http.Request, username, password string) (interface{}, bool) {
	expectedUser, expectedPass, _ := strings.Cut(h.CFG.Auth.CREDS, ":")
	user := user.User{}
	return user, allowUser(username, password, expectedUser, expectedPass, false)
}

func (h *Handler) NewBasicAuthFile() Middleware {
	return newBasicAuth(h.AllowUserFile)
}

func (h *Handler) AllowUserFile(r *http.Request, username, password string) (interface{}, bool) {
	return user.User{}, false
}

func (h *Handler) NewBasicAuthDB() Middleware {
	return newBasicAuth(h.AllowUserDB)
}

func (h *Handler) AllowUserDB(r *http.Request, username, password string) (interface{}, bool) {
	return user.User{}, false
}
