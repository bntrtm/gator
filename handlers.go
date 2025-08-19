package main

import (
	"errors"
	"fmt"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/bntrtm/gator/internal/database"
	"github.com/bntrtm/gator/internal/rss"
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

func handlerAgg(s *state, cmd command) error {
	rssFeed, err := rss.FetchFeed(s.client.httpClient, context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		err := errors.New("ERROR: Not enough arguments")
		return err
	}
	feedName := cmd.args[0]
	feedUrl := cmd.args[1]
	username := s.config.Username

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("%s added feed: %s", username, feed.Name))
	
	feedFollowParams := database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		err := errors.New("ERROR: Write feed url to follow")
		return err
	}
	url := cmd.args[0]
	username := s.config.Username

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:	   user.ID,
		FeedID:	   feed.ID,
	}
	s.db.CreateFeedFollow(context.Background(), params)
	fmt.Println(fmt.Sprintf("%s now following feed: %s", username, feed.Name))
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		err := errors.New("ERROR: Write feed url to unfollow")
		return err
	}
	url := cmd.args[0]

	params := database.DelFeedFollowParams{
		Name:	s.config.Username,
		Url:	url,
	}
	err := s.db.DelFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	username := s.config.Username

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Println(fmt.Sprintf("%s is not following any feeds.", username))
		return nil
	}
	fmt.Println(fmt.Sprintf("%s is following the following feeds:", username))
	for _, feed := range feeds {
		fmt.Println(fmt.Sprintf(feed.FeedName))
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
		feeds, err := s.db.ListFeeds(context.Background())
		if err != nil {
			return err
		}
		for i, feed := range feeds {
			fmt.Println(fmt.Sprintf("FEED #%d:", i + 1))
			output := fmt.Sprintf("%-12s %s\n", "\tName:", feed.Name)
			output += fmt.Sprintf("%-12s %s\n", "\tURL:", feed.Url)
			createdByUser, err := s.db.GetUserByID(context.Background(), feed.UserID)
			if err != nil {
				return err
			}
			output += fmt.Sprintf("%-12s %s\n", "\tCreated by:", createdByUser.Name)
			fmt.Println(output)
		}
		return nil
	}