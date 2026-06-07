package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hardiing/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	programState := state{
		cfg: &cfg,
	}
	commandsMap := make(map[string]func(*state, command) error)
	programCommands := commands{
		commands: commandsMap,
	}
	programCommands.register("login", handlerLogin)

	input := os.Args

	if len(input) < 2 {
		log.Fatal("need two arguments")
	}

	//programName := input[0]
	inputCommand := input[1]
	inputArguments := input[2:]

	if len(inputArguments) == 0 {
		fmt.Println("missing username argument")
		os.Exit(1)
	}

	userCommand := command{
		name: inputCommand,
		args: inputArguments,
	}

	programCommands.run(&programState, userCommand)
}
