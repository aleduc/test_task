package main

import (
	"flag"
	"log"

	"test_task/internal/app"
	"test_task/internal/config"
	"test_task/internal/config/loader"
)

func main() {
	var (
		cfg            config.Config
		configFilename = flag.String("configFilename", "app-config", "config file name")
		configDir      = flag.String("configDir", ".", "config directory")
	)
	err := loader.LoadConfig(*configDir, *configFilename, &cfg)
	if err != nil {
		log.Fatalf("load config: %s", err.Error())
	}
	app.Start(cfg)
}
