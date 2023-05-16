package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/it-novum/binaryd/config"
	"github.com/it-novum/binaryd/utils"
	"github.com/it-novum/binaryd/webserver"
)

// Daniel Ziegler it-novum GmbH 2023
// MIT License

func main() {
	// First try to parse the config file
	configPath := "./binaryd.ini"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	if !utils.FileExists(configPath) {
		fmt.Printf("Config file %v does not exists\n", configPath)
		os.Exit(1)
	}

	cfg := config.NewConfig(configPath)
	err := cfg.LoadIni()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cfg.ParseIni()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start Webserver
	ctx := context.Background()

	server := webserver.NewHttpServer(cfg)
	server.SetupHttpHandler()
	server.Start(ctx)

	// Setup signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown()

	// Wait for ListenAndServe goroutine to close.
	fmt.Println("Web server shutdown")

}
