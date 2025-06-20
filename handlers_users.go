package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BadgerCannon/boot-blog-agg/internal/database"
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
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      username,
		})
		if err != nil {
			return err
		}
		log.Printf("Registered %v\n", dbUser.Name)
		log.Println(dbUser)

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

func handlerListUsers(s *state, cmd command) error {

	if len(cmd.Args) > 0 {
		return fmt.Errorf("too many arguments")
	}
	allUsers, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	log.Println("Registered Users:")
	for _, user := range allUsers {
		if user.Name == s.config.CurrentUserName {
			log.Println("* ", user.Name, "(current)")
		} else {
			log.Println("* ", user.Name)
		}
	}
	return nil
}

func handlerResetDb(s *state, cmd command) error {

	if len(cmd.Args) > 0 {
		return fmt.Errorf("too many arguments")
	}
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}
	log.Println("All users deleted")
	return nil
}

func handlerFollowing(s *state, cmd command) error {

	expected_args := 0
	l := len(cmd.Args)
	switch {
	case l < expected_args || l > expected_args:
		return fmt.Errorf("incorrect number of arguments, expected %v got %v", expected_args, l)
	default:
		dbUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to lookup user in db: %w", err)
		}

		dbFeedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), dbUser.ID)
		if err != nil {
			return fmt.Errorf("failed to get follows for user in db: %w", err)
		}

		for _, follow := range dbFeedFollows {
			fmt.Println(follow)
		}

		return nil
	}
}
