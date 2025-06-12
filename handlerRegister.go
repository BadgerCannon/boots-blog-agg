package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/BadgerCannon/boots-go-blog-agg/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {

	switch len(cmd.Args) {
	case 0:
		return errors.New("no username provided")

	case 1:
		username := cmd.Args[0]
		log.Printf("Registering user %v\n", username)
		dbUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		})
		if err != nil {
			return err
		}
		log.Printf("Registered %v\n", dbUser.Name)
		log.Printf("full dbUser object: %v\n", dbUser)

		err = s.config.SetUser(dbUser.Name)
		if err != nil {
			return err
		}
		log.Printf("Logged in %v\n", dbUser.Name)
		return nil

	default:
		return errors.New("too many arguments provided, expected only username")
	}

}
