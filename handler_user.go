package main

import (
	"context"
	"fmt"

	"os"
	"strconv"
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
		return fmt.Errorf("reset unsuccessful: %w\n", err)
	}
	fmt.Println("users table reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	rows, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error in accessing table: %w\n", err)
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
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error parsing time duration: %v\n", err)
	}
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAdd(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Printf("Not enough arguments for %s\n", cmd.name)
		os.Exit(1)
		return nil
	}
	ctx := context.Background()

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
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("Error creating Feed Follow: %v\n", err)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	rows, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error in accessing feeds table: %w\n", err)
	}
	for _, row := range rows {
		fmt.Printf("Feed Name: %s\n", row.Name)
		fmt.Printf("Feed URL: %s\n", row.Url)
		fmt.Printf("Added by User: %s\n", row.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("Follow requires just one argument")
		os.Exit(1)
	}

	result, err := followFeed(s, user, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error running followFeed: %v\n", err)
	}
	fmt.Printf("Feed name: %s\n", result.FeedName)
	fmt.Printf("Current user: %s\n", result.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		fmt.Println("No arguments required")
		os.Exit(1)
	}
	ctx := context.Background()
	rows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("Error with GetFeedFollowsForUser query: %v\n", err)
	}
	fmt.Printf("Current User: %s\n", user.Name)
	for _, row := range rows {
		fmt.Printf("Feed: %s\n", row.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Unfollow requires just one argument.\n")
	}
	return unfollowFeed(s, user, cmd.args[0])

}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	ctx := context.Background()
	if len(cmd.args) != 1 {
		limit = 2
	} else {
		convertArg, err := strconv.ParseInt(cmd.args[0], 10, 32)
		limit = int32(convertArg)
		if err != nil {
			fmt.Printf("error converting argument to int")
		}
	}
	browseFeeds, err := s.db.GetPostsForUser(ctx, limit)
	if err != nil {
		return fmt.Errorf("Error getting posts for user: %v\n", err)
	}
	for _, post := range browseFeeds {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description)
		fmt.Printf("Published At: %s\n", post.PublishedAt.Time)
	}
	return nil
}
