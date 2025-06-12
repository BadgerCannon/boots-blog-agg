package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/BadgerCannon/boots-go-blog-agg/internal/config"

	"github.com/BadgerCannon/boots-go-blog-agg/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("ERROR: Not enough arguments provided")
	}

	activeConfig, err := config.Read()
	if err != nil {
		log.Fatal("ERROR: Failed to load config: ", "err", err)
	}

	db, err := sql.Open("postgres", activeConfig.DbUrl)
	if err != nil {
		log.Fatal("ERROR: Unable to connect to database", "err", err)
	}

	activeState := state{
		config: &activeConfig,
		db:     database.New(db),
	}

	availableCommands := commands{
		commands: make(map[string]func(*state, command) error),
	}

	availableCommands.register("login", handlerLogin)
	availableCommands.register("register", handlerRegister)

	// slog.Debug("msg", "activeState", activeState, "os.Args", os.Args, "availableCommands", availableCommands)

	err = availableCommands.run(&activeState, command{
		Name: os.Args[1],
		Args: os.Args[2:],
	})
	if err != nil {
		log.Fatalf("ERROR: command failed: %v\n", err)
	}
}
