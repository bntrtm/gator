package main

import (
	"errors"
	"fmt"

	"github.com/bntrtm/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		err := errors.New("ERROR: Username required")
		return err
	}
	s.config.SetUser(cmd.args[0])
	fmt.Println(fmt.Sprintf("Logged in as '%s'", s.config.Username))
	return nil
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

