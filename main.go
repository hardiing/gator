package main

import (
	"github.com/hardiing/gator/internal/config"
)

func main() {
	cfg := config.Read()
	config.SetUser(cfg)
	config.Read()
}
