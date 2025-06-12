package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	switch len(cmd.Args) {
	case 0:
		return errors.New("no username provided")
	case 1:
		username := cmd.Args[0]
		_, err := s.db.GetUser(context.Background(), username)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no registered user '%v'", username)
		} else if err != nil {
			return fmt.Errorf("failed to get user from database: %w", err)
		}
		err = s.config.SetUser(username)
		if err != nil {
			return err
		}
		log.Printf("Logged in %v\n", username)
		return nil
	default:
		return errors.New("too many arguments provided, expected only username")
	}
}
