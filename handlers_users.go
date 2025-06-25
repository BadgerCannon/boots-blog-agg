package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/BadgerCannon/boot-blog-agg/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	usage := "Usage: boot-blog-agg feeds USERNAME\nRegister a new user and log them in."
	args_len := len(cmd.Args)
	checkUsage(1, 1, args_len, usage)

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

}

func handlerLogin(s *state, cmd command) error {
	usage := "Usage: boot-blog-agg login USERNAME\nLogin as the given user."
	args_len := len(cmd.Args)
	checkUsage(1, 1, args_len, usage)

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

}

func handlerListUsers(s *state, cmd command) error {
	usage := "Usage: boot-blog-agg users\nList configured users."
	args_len := len(cmd.Args)
	checkUsage(0, 0, args_len, usage)

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
	usage := "Usage: boot-blog-agg reset\nDelete all configured users. INTENDED FOR DEV USE ONLY."
	args_len := len(cmd.Args)
	checkUsage(0, 0, args_len, usage)

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}
	log.Println("All users deleted")
	return nil
}

// Logged in Handlers
func handlerFollowing(s *state, cmd command, user database.User) error {
	usage := "Usage: boot-blog-agg following\nPrint the list of feeds the current user is following to the terminal"
	args_len := len(cmd.Args)
	checkUsage(0, 0, args_len, usage)

	dbFeedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get follows for user in db: %w", err)
	}

	for _, follow := range dbFeedFollows {
		fmt.Println(follow)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	usage := "Usage: boot-blog-agg browse [LIMIT] \nPrint LIMIT (default 2) most recent scraped posts to the terminal"
	args_len := len(cmd.Args)
	checkUsage(0, 1, args_len, usage)

	var postLimit int32 // database.GetPostsforUserParams requires int32
	if args_len == 1 {
		var err error
		limitArg, err := strconv.ParseInt(cmd.Args[0], 10, 32) // returns int64 _limited_ to int32 range
		if err != nil {
			return fmt.Errorf("post limit must be numeric")
		}
		postLimit = int32(limitArg)
	} else {
		postLimit = 2
	}

	dbPosts, err := s.db.GetPostsforUser(context.Background(), database.GetPostsforUserParams{
		UserID: user.ID,
		Limit:  postLimit,
	})
	if err != nil {
		return fmt.Errorf("failed to get posts for user in db: %w", err)
	}

	for _, post := range dbPosts {
		fmt.Println("==========================================================================================================")
		fmt.Println(post)
	}
	fmt.Println("==========================================================================================================")

	return nil
}
