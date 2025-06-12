package main

import (
	"context"
	"fmt"
	"log"
)

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
