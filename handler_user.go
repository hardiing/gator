package main

import (
	"context"
	"fmt"
	"html"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hardiing/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no arguments found")
	}

	err := s.cfg.SetUser(cmd.args[0])
	_, err = s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("user %s not found\n", cmd.args[0])
		os.Exit(1)
		return err
	}

	fmt.Printf("User %s has been set\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no arguments found")
	}

	ctx := context.Background()
	newID := uuid.New()
	time := time.Now()
	params := database.CreateUserParams{
		ID:        newID,
		CreatedAt: time,
		UpdatedAt: time,
		Name:      cmd.args[0],
	}

	result, err := s.db.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("username %s already exists\n", params.Name)
		os.Exit(1)
	}

	s.cfg.SetUser(result.Name)
	fmt.Printf("user %s created\n", result.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		fmt.Errorf("reset unsuccessful: %w\n", err)
		os.Exit(1)
	}
	fmt.Println("users table reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	rows, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Errorf("error in accessing table: %w\n", err)
		os.Exit(1)
	}
	for _, row := range rows {
		if row == s.cfg.Username {
			fmt.Printf("* %s (current)\n", row)
		} else {
			fmt.Printf("* %s\n", row)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		fmt.Printf("Error found fetching feed: %v\n", err)
	}
	fmt.Printf("Title: %s\n", html.UnescapeString(feed.Channel.Title))
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", html.UnescapeString(feed.Channel.Description))
	for _, item := range feed.Channel.Item {
		fmt.Printf("Item Title: %s\n", html.UnescapeString(item.Title))
		fmt.Printf("Item Link: %s\n", item.Link)
		fmt.Printf("Item Description: %s\n", html.UnescapeString(item.Description))
		fmt.Printf("Item PubDate: %s\n", item.PubDate)
	}
	return nil
}

func handlerAdd(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		fmt.Printf("Not enough arguments for %s\n", cmd.name)
		os.Exit(1)
		return nil
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.Username)
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		os.Exit(1)
	}
	feed, err := addFeed(s, user, cmd.args[0], cmd.args[1])
	if err != nil {
		fmt.Printf("Add feed failed: %v\n", err)
		os.Exit(1)
		return err
	}
	fmt.Println("Adding to feeds table:")
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created at: %v\n", feed.CreatedAt)
	fmt.Printf("Updated at: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("User ID: %v\n", feed.UserID)
	return nil
}
