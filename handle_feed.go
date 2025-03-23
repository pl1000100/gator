package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/gator/internal/database"
)

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
	err = scrapeFeed(s, nextFeed)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func scrapeFeed(s *state, nextFeed database.Feed) error {
	markParams := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err := s.db.MarkFeedFetched(context.Background(), markParams)
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
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{Valid: false},
			Url:         item.Link,
			Description: sql.NullString{Valid: false},
			PublishedAt: sql.NullTime{Valid: false},
			FeedID:      nextFeed.ID,
		}
		if len(item.Title) > 0 {
			postParams.Title = sql.NullString{String: item.Title, Valid: true}
		}
		if len(item.Description) > 0 {
			postParams.Description = sql.NullString{String: item.Description, Valid: true}
		}
		if len(item.PubDate) > 0 {
			t, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				fmt.Println("Error parsing RFC1123Z:", err)
				return err
			}
			postParams.PublishedAt = sql.NullTime{Time: t, Valid: true}
		}

		post, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Println(err)
			return err
		}
		fmt.Println(post.Title.String, "succesfully saved.")
	}

	return nil
}
