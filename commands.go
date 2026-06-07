package main

import (
	"fmt"

	"github.com/hardiing/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	input, ok := c.commands[cmd.name]
	if ok {
		return input(s, cmd)
	}
	return fmt.Errorf("command %s not found", cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
