package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"

	"github.com/BadgerCannon/boots-go-blog-agg/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}

	var feed RSSFeed
	xd := xml.NewDecoder(resp.Body)
	err = xd.Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode xml: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	fmt.Println(feed)

	return &feed, nil
}

func handlerAgg(s *state, cmd command) error {
	fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	l := len(cmd.Args)
	switch {
	case l < 2 || l > 2:
		return fmt.Errorf("incorrect number of arguments, expected 2 got %v", l)

	case l == 2:
		if _, err := url.ParseRequestURI(cmd.Args[1]); err != nil {
			return fmt.Errorf("failed to parse URL '%v'", cmd.Args[1])
		}

		dbUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to lookup user in db: %w", err)
		}

		dbFeed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
			Name:   cmd.Args[0],
			Url:    cmd.Args[1],
			UserID: dbUser.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add feed to db: %w", err)
		}
		log.Printf("dbFeed: %v\n", dbFeed)
		return nil

	default:
		return fmt.Errorf("unexpected (impossible?) switch fallthrough")
	}
}
