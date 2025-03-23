package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/gator/internal/database"
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
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator/1.0")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()
	data := RSSFeed{}
	if err := xml.Unmarshal(body, &data); err != nil {
		return &RSSFeed{}, err
	}
	return &data, nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		fmt.Printf("error: wrong number of arguments passed, expected 1, got %d", len(cmd.Arguments))
		os.Exit(1)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		fmt.Printf("error: wrong number of arguments passed, expected 2, got %d", len(cmd.Arguments))
		os.Exit(1)
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	newCmd := command{
		Name:      "follow",
		Arguments: cmd.Arguments[1:],
	}
	handlerFollow(s, newCmd, user)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		fmt.Printf("error: wrong number of arguments passed, expected 0, got %d", len(cmd.Arguments))
		os.Exit(1)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feeds {
		fmt.Println(f.Name, f.Url, f.Username)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	params := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}

	return nil
}
