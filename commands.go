package main

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (commands *commands) run(s *state, cmd command) error {
	if cmd_func, ok := commands.commands[cmd.Name]; ok {
		return cmd_func(s, cmd)
	} else {
		return fmt.Errorf("no command called '%v' registered.\n\n%v", cmd.Name, commands.listCommands())
	}
}

func (commands *commands) register(name string, f func(*state, command) error) error {
	if _, ok := commands.commands[name]; !ok {
		commands.commands[name] = f
	} else {
		return fmt.Errorf("command with Name '%v' already registered", name)
	}
	return nil
}

func checkUsage(min_args, max_args, arg_count int, usage string) {
	if arg_count < min_args {
		log.Fatalf("too few arguments, expected %v got %v\n\n%v\n", min_args, arg_count, usage)
	} else if arg_count > max_args {
		log.Fatalf("too many arguments, expected %v got %v\n\n%v\n", max_args, arg_count, usage)
	}
}

func (commands *commands) listCommands() string {
	commandNames := slices.Sorted(maps.Keys(commands.commands))

	availableCommands := strings.Join(commandNames, ", ")

	return fmt.Sprintf("Available Commands:\n\t%v", availableCommands)
}
