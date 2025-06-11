package main

import (
	"context"
	"internal/config"
	"log/slog"
	"os"
)

type state struct {
	config *config.Config
}

func main() {
	slog.Default().Enabled(context.TODO(), slog.LevelDebug)
	// var programLevel = new(slog.LevelVar)
	// h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	// slog.SetDefault(slog.New(h))
	// programLevel.Set(slog.LevelDebug)

	if len(os.Args) < 2 {
		slog.Error("Not enough arguments provided")
		os.Exit(1)
	}

	activeConfig, err := config.Read()
	if err != nil {
		slog.Error("FATAL: Failed to load config: ", "err", err)
		os.Exit(1)
	}

	activeState := state{
		config: &activeConfig,
	}
	availableCommands := commands{
		commands: make(map[string]func(*state, command) error),
	}

	availableCommands.register("login", handlerLogin)
	slog.Debug("msg", "activeState", activeState, "os.Args", os.Args, "availableCommands", availableCommands)

	err = availableCommands.run(&activeState, command{
		Name: os.Args[1],
		Args: os.Args[2:],
	})
	if err != nil {
		slog.Error("Failed to run command", "err", err)
		os.Exit(1)
	}
}
