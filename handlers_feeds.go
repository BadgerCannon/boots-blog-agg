package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/BadgerCannon/boot-blog-agg/internal/database"
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

func (feed RSSFeed) String() string {
	// return fmt.Sprintf("Feed Title: 			%v\nFeed Description: 		%v\nFeed Link: 			%v\nFeed Post Count:		%v", feed.Channel.Title, feed.Channel.Description, feed.Channel.Link, len(feed.Channel.Item))
	return fmt.Sprintf("Feed Title: 			%v\nFeed Link: 			%v\nFeed Post Count:		%v", feed.Channel.Title, feed.Channel.Link, len(feed.Channel.Item))
}

func (item RSSItem) String() string {
	// return fmt.Sprintf("| Post Title: 			%v\n| Post Description: 		%v\n| Post Date: 			%v\n| Post Link: 			%v", item.Title, item.Description, item.PubDate, item.Link)
	return fmt.Sprintf("| Post Title: 			%v\n| Post Date: 			%v\n| Post Link: 			%v", item.Title, item.PubDate, item.Link)
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

	return &feed, nil
}

func handlerListFeeds(s *state, cmd command) error {

	dbFeeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds from db: %w", err)
	}

	for _, feed := range dbFeeds {
		dbUser, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user from db: %w", err)
		}
		log.Println("Feed created by", dbUser.Name)
		log.Println(feed)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	expected_args := 1
	l := len(cmd.Args)
	switch {
	case l < expected_args || l > expected_args:
		return fmt.Errorf("incorrect number of arguments, expected %v got %v", expected_args, l)

	default:
		timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("failed to parse user defined duration : %w", err)
		}
		ticker := time.NewTicker(timeBetweenRequests)
		for ; ; <-ticker.C {
			scrapeFeeds(s)
		}
	}
}

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("ERROR: failed to get feed to fetch from db: %v", err)
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Printf("ERROR: failed to fetch feed: %v", err)
	}
	s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	fmt.Println("+======================================================================+")
	fmt.Println(feed)
	fmt.Println("+----------------------------------------------------------------------+")
	for _, post := range feed.Channel.Item {
		fmt.Println(post)
		fmt.Println("+----------------------------------------------------------------------+")
	}
}

// Logged in Handlers

func handlerAddFeed(s *state, cmd command, user database.User) error {
	l := len(cmd.Args)
	switch {
	case l < 2 || l > 2:
		return fmt.Errorf("incorrect number of arguments, expected 2 got %v", l)

	case l == 2:
		if _, err := url.ParseRequestURI(cmd.Args[1]); err != nil {
			return fmt.Errorf("failed to parse URL '%v'", cmd.Args[1])
		}

		dbFeed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
			Name:   cmd.Args[0],
			Url:    cmd.Args[1],
			UserID: user.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add feed to db: %w", err)
		}
		log.Println(dbFeed)

		dbFeedFollow, err := s.db.AddFeedFollow(context.Background(), database.AddFeedFollowParams{
			UserID: user.ID,
			FeedID: dbFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add feed follow to db: %w", err)
		}
		log.Println(dbFeedFollow)

		return nil

	default:
		return fmt.Errorf("unexpected (impossible?) switch fallthrough")
	}
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	expected_args := 1
	l := len(cmd.Args)
	switch {
	case l < expected_args || l > expected_args:
		return fmt.Errorf("incorrect number of arguments, expected %v got %v", expected_args, l)

	default:
		if _, err := url.ParseRequestURI(cmd.Args[0]); err != nil {
			return fmt.Errorf("failed to parse URL '%v'", cmd.Args[1])
		}

		dbFeed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
		if err != nil {
			return fmt.Errorf("failed to lookup feed in db: %w", err)
		}

		dbFeedFollow, err := s.db.AddFeedFollow(context.Background(), database.AddFeedFollowParams{
			UserID: user.ID,
			FeedID: dbFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add feed follow to db: %w", err)
		}
		log.Println(dbFeedFollow)
		return nil
	}
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	expected_args := 1
	l := len(cmd.Args)
	switch {
	case l < expected_args || l > expected_args:
		return fmt.Errorf("incorrect number of arguments, expected %v got %v", expected_args, l)

	default:
		if _, err := url.ParseRequestURI(cmd.Args[0]); err != nil {
			return fmt.Errorf("failed to parse URL '%v'", cmd.Args[1])
		}

		dbFeed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
		if err != nil {
			return fmt.Errorf("failed to lookup feed in db: %w", err)
		}

		err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: dbFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to unfollow feed: %w", err)
		}
		return nil
	}
}
