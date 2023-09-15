// Package main is the main package
package main

import (
	"flag"

	"github.com/SerjRamone/not-found-bot/config"
	"github.com/SerjRamone/not-found-bot/internal/app"
)

func main() {
	configFile := flag.String("p", ".env", "Path to config file")
	flag.Parse()

	cfg := config.Get(*configFile)

	app.Run(cfg)
}
