package main

import (
	"fmt"
	"context"
	"errors"
	"github.com/bntrtm/gator/internal/database"
)

func middlewareLoggedIn(
	handler func(
		s *state, 
		cmd command, 
		user database.User,
		) error,
) (func(*state, command) error) {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.Username)
		if err != nil {
			err = errors.New(fmt.Sprintf("ERROR: User '%s' is not logged in", s.config.Username))
			return err
		}

		return handler(s, cmd, user)
	}
}