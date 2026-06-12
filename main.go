package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/hardiing/gator/internal/config"
	"github.com/hardiing/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Url)
	dbQueries := database.New(db)

	programState := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	commandsMap := make(map[string]func(*state, command) error)
	programCommands := commands{
		commands: commandsMap,
	}
	programCommands.register("login", handlerLogin)
	programCommands.register("register", handlerRegister)
	programCommands.register("reset", handlerReset)

	input := os.Args

	if len(input) < 2 {
		log.Fatal("need two arguments")
	}

	//programName := input[0]
	inputCommand := input[1]
	inputArguments := input[2:]

	userCommand := command{
		name: inputCommand,
		args: inputArguments,
	}

	programCommands.run(&programState, userCommand)
}
