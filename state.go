package main

import (
	"errors"
	"fmt"
	"github.com/bntrtm/gator/internal/config"
	"github.com/bntrtm/gator/internal/database"
)

type state struct {
	db *database.Queries
	config *config.Config
	client *Client
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New(fmt.Sprintf("Unknown command '%s'", cmd.args[0]))
	}
	return fn(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.handlers[name]
	if ok {
		fmt.Println(fmt.Sprintf("ERROR: Command '%s' already exists in registry", name))
	}
	c.handlers[name] = f
	return
}

