package rss

import (
	"io"
	"fmt"
	"html"
	"context"
	"net/http"
	"encoding/xml"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// DEPRECATED
func(r *RSSFeed) PrintItems() {
	for _, item := range r.Channel.Item {
		fmt.Println(item.Title)
	}
	return
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func FetchFeed(client http.Client, ctx context.Context, feedURL string) (*RSSFeed, error) {
	
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(xmlData, &rssFeed)
	if err != nil {
		return nil, err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for _, rssItem := range rssFeed.Channel.Item {
		rssItem.Title = html.UnescapeString(rssItem.Title)
		rssItem.Description = html.UnescapeString(rssItem.Description)
	}

	return &rssFeed, nil
}