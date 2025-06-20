package main

import (
	"context"
	"fmt"

	"github.com/BadgerCannon/boot-blog-agg/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		dbUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to lookup user in db: %w", err)
		}
		return handler(s, cmd, dbUser)
	}
}
