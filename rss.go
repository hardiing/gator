package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hardiing/gator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error: %v\n", err)
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

func unfollowFeed(s *state, user database.User, url string) error {
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("Error getting feed by url: %v\n", err)
	}
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("Error deleting feed follow: %v\n", err)
	}
	return err
}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	next, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("Error getting next feed to fetch: %v\n", err)
		return
	}
	_, err = s.db.MarkFeedFetched(ctx, next.ID)
	if err != nil {
		fmt.Printf("Error marking as fetched: %v\n", err)
		return
	}
	fetch, err := fetchFeed(ctx, next.Url)
	if err != nil {
		fmt.Printf("Error getting feed by URL: %v\n", err)
		return
	}
	for _, item := range fetch.Channel.Item {
		pubDateTime, err := time.Parse(time.RFC3339, item.PubDate)
		if err != nil {
			pubDateTime, err = time.Parse(time.RFC1123Z, item.PubDate)
		}

		pt := sql.NullTime{Valid: err == nil}
		if err == nil {
			pt.Time = pubDateTime
		}
		now := time.Now()

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         next.Url,
			Description: item.Description,
			PublishedAt: pt,
			FeedID:      next.ID,
		}
		_, err = s.db.CreatePost(ctx, params)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				// "23505" represents unique_violation in Postgres
				if pqErr.Code == "23505" {
					continue
				} else {
					fmt.Printf("Error encountered: %s\n", err)
				}
			} else {
				fmt.Printf("Error creating post: %v\n", err)
			}
		}
	}
}
