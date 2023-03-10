package main

import (
	"context"
	"flag"
	"os"

	"github.com/devops-rob/boundary-census/config"
	"github.com/devops-rob/boundary-census/handlers"

	bc "github.com/devops-rob/boundary-census/clients/boundary"
	nc "github.com/devops-rob/boundary-census/clients/nomad"

	"github.com/hashicorp/go-hclog"
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

	boundaryClient, err := bc.New(
		cfg.Boundary.Address,
		cfg.Boundary.OrgID,
		cfg.Boundary.DefaultProject,
		cfg.Boundary.AuthMethodID,
		map[string]interface{}{
			"login_name": cfg.Boundary.Username,
			"password":   cfg.Boundary.Password,
		},
		cfg.Boundary.Enterprise,
	)
	if err != nil {
		logger.Error("Unable to create client", "error", err)
		os.Exit(1)
	}

	// Create a context to cancel the event subscription when the program exits
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create the handler
	targetHandler := handlers.NewTarget(logger, boundaryClient)

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

			switch t {
			case "AllocationUpdated":
				alloc, err := e.Allocation()
				if err != nil {
					logger.Error("unable to fetch allocation", "error", err)
				}

				logger.Info("handle allocation", "status", alloc.ClientStatus)
				switch alloc.ClientStatus {
				case "running":
					for _, n := range alloc.AllocatedResources.Shared.Networks {
						ports := []uint32{}

						for _, p := range n.DynamicPorts {
							ports = append(ports, uint32(p.Value))
						}

						si := &handlers.ServiceInstance{
							Location: n.IP,
							Ports:    ports,
						}

						// call create
						ids, err := targetHandler.Create(
							si,
							alloc.Name,
							cfg.Boundary.OrgID,
							cfg.Boundary.DefaultProject,
							cfg.Boundary.DefaultIngressFilter,
							cfg.Boundary.DefaultEgressFilter,
						)

						if err != nil {
							logger.Error("Unable to create tasks", "error", err)
							break
						}

						logger.Info("Created tasks", "ids", ids)
					}
				case "complete":
					err := targetHandler.DeleteWithPrefix(alloc.Name, cfg.Boundary.OrgID, cfg.Boundary.DefaultProject)
					if err != nil {
						logger.Error("Unable to delete target", "prefix", alloc.Name, "error", err)
					}

				}
			}
		}
	}
}
