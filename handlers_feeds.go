package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/BadgerCannon/boot-blog-agg/internal/database"
	"github.com/lib/pq"
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
	usage := "Usage: boot-blog-agg feeds\nList the feeds configured for scraping, and which user"
	args_len := len(cmd.Args)
	checkUsage(0, 0, args_len, usage)

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
	usage := "Usage: boot-blog-agg register USERNAME\nRegister a new user and login."
	args_len := len(cmd.Args)
	checkUsage(0, 0, args_len, usage)

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to parse user defined duration : %w", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func parsePostDate(post_date string) (time.Time, error) {
	layouts := []string{
		time.RFC822,
		time.RFC1123Z, // matches examples
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC3339,
		time.RFC3339Nano,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
	}
	for _, layout := range layouts {
		parsedDate, err := time.Parse(layout, post_date)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date '%v' with configured layouts", post_date)
}

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	fmt.Printf("INFO: Scraping %v\n", nextFeed.Name)
	if err != nil {
		fmt.Printf("ERROR: failed to get feed to fetch from db: %v", err)
		return
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Printf("ERROR: failed to fetch feed: %v", err)
		return
	}
	s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	for _, post := range feed.Channel.Item {
		postDate, err := parsePostDate(post.PubDate)
		if err != nil {
			fmt.Printf("WARN: failed to save post %v to db: %v\n", post.Link, err)
			continue
		}
		_, err = s.db.AddPost(context.Background(), database.AddPostParams{
			FeedID:      nextFeed.ID,
			Title:       post.Title,
			Description: sql.NullString{String: post.Description, Valid: true},
			Url:         post.Link,
			PublishedAt: postDate,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code != "23505" {
					fmt.Printf("WARN: failed to save post %v to db: %v\n", post.Link, err)
					continue
				}
			} else {
				fmt.Printf("err: %v\n", err)
				continue
			}
		}
	}
}

// Logged in Handlers

func handlerAddFeed(s *state, cmd command, user database.User) error {
	usage := "Usage: boot-blog-agg addfeed NAME URL\nAdd a feed to gator for scraping."
	args_len := len(cmd.Args)
	checkUsage(2, 2, args_len, usage)

	if _, err := url.ParseRequestURI(cmd.Args[1]); err != nil { // don't need a url.*URL struct, just need to check the string is a valid url
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
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	usage := "Usage: boot-blog-agg follow URL\nFollow a feed that has already been added to gator by another user by URL."
	args_len := len(cmd.Args)
	checkUsage(1, 1, args_len, usage)

	if _, err := url.ParseRequestURI(cmd.Args[0]); err != nil { // don't need a url.*URL struct, just need to check the string is a valid url
		return fmt.Errorf("failed to parse URL '%v'", cmd.Args[0])
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

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	usage := "Usage: boot-blog-agg unfollow URL\nUnfollow a feed by URL."
	args_len := len(cmd.Args)
	checkUsage(1, 1, args_len, usage)

	if _, err := url.ParseRequestURI(cmd.Args[0]); err != nil {
		return fmt.Errorf("failed to parse URL '%v'", cmd.Args[0])
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
