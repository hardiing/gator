package config

import (
	"fmt"
	"os"
)

const configFileName = " .gatorconfig.json"

type Config struct {
	url      string `json:"db_url"`
	username string `json:"current_user_name"`
}

func Read() Config {
	// Export a Read function that reads the JSON file found at ~/.gatorconfig.json
	// and returns a Config struct. It should read the file from the HOME directory,
	// then decode the JSON string into a new Config struct. I used os.UserHomeDir to get the location of HOME.
	home := os.UserHomeDir
	fmt.Printf("%s", home)
}

func SetUser(cfg Config) Config {
	// Export a SetUser method on the Config struct that writes the config struct
	// to the JSON file after setting the current_user_name field.
}

func getConfigFilePath() (string, error) {

}

func write(cfg Config) error {

}
