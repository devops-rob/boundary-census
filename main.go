package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/devops-rob/boundary-census/config"

	nc "github.com/devops-rob/boundary-census/clients/nomad"

	"github.com/hashicorp/go-hclog"
	"github.com/kr/pretty"
)

var (
	configFile = flag.String("config", "./config.hcl", "path to the Census config file")
	logLevel   = flag.String("log_level", "info", "log level, info, debug, trace")
)

func main() {
	opts := hclog.DefaultOptions
	opts.Color = hclog.AutoColor
	opts.Level = hclog.LevelFromString(*logLevel)

	logger := hclog.New(opts)

	// load the config
	flag.Parse()

	cfg, err := config.Parse(*configFile)
	if err != nil {
		logger.Error("Unable to read config file", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting event stream", "addr", cfg.Nomad.Address)

	streamConfig := nc.DefaultClientConfig()
	streamConfig.Address = cfg.Nomad.Address

	// Create a new Stream object
	s := nc.NewStream(&streamConfig)

	// Create a context to cancel the event subscription when the program exits
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to deployment events
	eventStream, err := s.Subscribe(ctx)
	if err != nil {
		s.L.Error("error subscribing to events", "error", err)
		os.Exit(1)
	}

	// Loop over the event stream channel to read events as they happen
	for event := range eventStream {
		for _, e := range event.Events {
			t := e.Type
			fmt.Println(t)

			switch t {
			case "AllocationUpdated":
				alloc, err := e.Allocation()
				if err != nil {
					logger.Error("unable to fetch allocation", "error", err)
				}

				pretty.Println(alloc)
			}
		}
	}
}
