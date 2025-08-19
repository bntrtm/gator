package main

import (
	"errors"
	"fmt"
	"time"
	"context"
	"strings"
	"database/sql"

	"github.com/google/uuid"
	"github.com/bntrtm/gator/internal/config"
	"github.com/bntrtm/gator/internal/database"
	"github.com/bntrtm/gator/internal/rss"
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
		return errors.New(fmt.Sprintf("Unknown command '%s'", cmd.name))
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

func ScrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return err
	}
	feed, err := rss.FetchFeed(s.client.httpClient, context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}
	saveItems(s, feed, feedToFetch.ID)
	return nil
}

func saveItems(s *state, r *rss.RSSFeed, feedID uuid.UUID) error {
	for _, item := range r.Channel.Item {

		fixedPubDate := sql.NullTime{}
		var parseErr error
		layouts := []string{
			time.RFC1123Z,
			time.RFC1123,
			time.RFC3339,
			time.RFC822,
		}
		for _, layout := range layouts {
			t, parseErr := time.Parse(layout, item.PubDate)
			if parseErr == nil {
				fixedPubDate = sql.NullTime{
					Time: t,
					Valid: true,
				}
				break
			}
		}
		if parseErr != nil {
			fmt.Println("ERROR: could not resolve layout for publish date: ", item.PubDate)
		}
		
		nullDesc := sql.NullString{
			String: item.Description,
			Valid: item.Description != "",
		}

		params := database.CreatePostParams{
			ID:				uuid.New(),
			CreatedAt:		time.Now(),
			UpdatedAt:		time.Now(),
			Title:			item.Title,
			Url:			item.Link,
			Description:	nullDesc,
			PublishedAt:	fixedPubDate,
			FeedID:			feedID,
		}
		_, err := s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Println("Couldn't create post: %v", err)
			continue
		}
		
	}

	return nil
}