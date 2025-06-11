package main

import (
	"errors"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	switch len(cmd.Args) {
	case 0:
		return errors.New("no username provided")
	case 1:
		username := cmd.Args[0]
		err := s.config.SetUser(username)
		if err != nil {
			return err
		}
		log.Printf("Username set to %v\n", username)
		return nil
	default:
		return errors.New("too many arguments provided, expected only username")
	}
}
