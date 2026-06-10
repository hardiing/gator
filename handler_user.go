package main

import (
	"context"
	"fmt"
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
	fmt.Printf("user %s created, %v\n", result.Name, result)
	return nil
}
