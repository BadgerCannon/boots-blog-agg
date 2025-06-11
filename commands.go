package main

import "fmt"

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
		return fmt.Errorf("no command called '%v' registered", cmd.Name)
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
