package store

import (
	"fmt"

	"github.com/vinser/flibgolite/internal/user"
	"golang.org/x/crypto/bcrypt"
)

// Auth
var ErrorBadUserOrPassword = fmt.Errorf("Bad user name or password")

func (db *DB) GetUserByUsername(username string) (user.User, error) {
	// TODO Replace with real request
	users := map[string]user.User{
		"admin": {
			ID:       0,
			Username: "admin",
			Password: "admin",
			Email:    "admin@localhost",
		},
		"john": {
			ID:       1,
			Username: "john",
			Password: "p@ss",
			Email:    "john@localhost",
		},
	}

	if u, ok := users[username]; ok {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
		if err != nil {
			return user.User{}, err
		}
		u.Password = string(hashed)
		return u, nil
	} else {
		return user.User{}, ErrorBadUserOrPassword
	}
}
