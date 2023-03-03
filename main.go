package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/devops-rob/boundary-census/config"

	bc "github.com/devops-rob/boundary-census/clients/boundary"
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

	// Create a new Stream object
	streamConfig := nc.DefaultClientConfig()
	streamConfig.Address = cfg.Nomad.Address
	s := nc.NewStream(&streamConfig)

	logger.Info("Creating Boundary client", "addr", cfg.Boundary.Address, "org", cfg.Boundary.OrgID, "auth_id", cfg.Boundary.AuthMethodID, "username", cfg.Boundary.Username)

	_, err = bc.New(
		cfg.Boundary.Address,
		cfg.Boundary.OrgID,
		cfg.Boundary.DefaultProject,
		cfg.Boundary.AuthMethodID,
		map[string]interface{}{
			"login_name": cfg.Boundary.Username,
			"password":   cfg.Boundary.Password,
		},
	)
	if err != nil {
		logger.Error("Unable to create client", "error", err)
		os.Exit(1)
	}

	// Create a context to cancel the event subscription when the program exits
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to deployment events
	logger.Info("Starting event stream", "addr", cfg.Nomad.Address)
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

				// create a service from the allocation
				//si := handlers.ServiceInstance{
				//}

				pretty.Println(alloc)
			}
		}
	}
}
