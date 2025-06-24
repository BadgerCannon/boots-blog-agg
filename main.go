package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/BadgerCannon/boot-blog-agg/internal/config"

	"github.com/BadgerCannon/boot-blog-agg/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	var cmdName string
	var cmdArgs []string
	if len(os.Args) == 1 {
		log.Fatal("ERROR: No command provided")
	} else if len(os.Args) == 2 {
		cmdName = os.Args[1]
		cmdArgs = []string{}
	} else {
		cmdName = os.Args[1]
		cmdArgs = os.Args[2:]
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
	availableCommands.register("reset", handlerResetDb)
	availableCommands.register("users", handlerListUsers)
	availableCommands.register("following", middlewareLoggedIn(handlerFollowing))
	availableCommands.register("browse", middlewareLoggedIn(handlerBrowse))

	availableCommands.register("agg", handlerAgg)
	availableCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	availableCommands.register("feeds", handlerListFeeds)
	availableCommands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	availableCommands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

	// slog.Debug("msg", "activeState", activeState, "os.Args", os.Args, "availableCommands", availableCommands)

	err = availableCommands.run(&activeState, command{
		Name: cmdName,
		Args: cmdArgs,
	})
	if err != nil {
		log.Fatalf("ERROR: command failed: %v\n", err)
	}
}
