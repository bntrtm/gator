package main

import (
	"errors"
	"fmt"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/bntrtm/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		err := errors.New("ERROR: Username required")
		return err
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		err = errors.New(fmt.Sprintf("ERROR: User '%s' does not exist", username))
		return err
	}
	s.config.SetUser(username)
	fmt.Println(fmt.Sprintf("Logged in as '%s'", s.config.Username))
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		err := errors.New("ERROR: Name is required")
		return err
	}
	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
	}
	userData, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	s.config.SetUser(cmd.args[0])
	fmt.Println(fmt.Sprintf("User '%s' created. Data: %s", s.config.Username, userData))
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		output := "* " + user.Name
		if s.config.Username == user.Name {
			output += " (current)"
		}
		fmt.Println(output)
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DelUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}