package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
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
	//fmt.Printf("Item Link: %s\n", feed.Channel.Item.Link)
	//fmt.Printf("Item Description: %s\n", feed.Channel.Item)
	return &feed, nil
}
