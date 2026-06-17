package main

import (
	"context"
	"fmt"

	"github.com/hardiing/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Username)
		if err != nil {
			return fmt.Errorf("Cannot find user: %s\n", user)
		}
		return handler(s, cmd, user)
	}
}
