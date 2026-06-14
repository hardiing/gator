package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hardiing/gator/internal/database"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		fmt.Printf("Error: %w\n", err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err := xml.Unmarshal(body, &feed); err != nil {
		fmt.Printf("Error encountered: %v\n", err)
	}
	return &feed, nil
}

func addFeed(s *state, user database.User, name, url string) (database.Feed, error) {
	ctx := context.Background()
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	return s.db.CreateFeed(ctx, params)
}

func followFeed(s *state, user database.User, url string) (database.CreateFeedFollowRow, error) {
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		fmt.Printf("GetFeedByURL failed: %v\n", err)
		os.Exit(1)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	return s.db.CreateFeedFollow(ctx, params)
}
