package stream

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
)

type Stream struct {
	nomad *api.Client
	L     hclog.Logger
}

type AllocationUpdate struct {
	Allocation *api.Allocation
	Job        *api.Job
	Deployment *api.Deployment
}

type ClientConfig struct {
	Address   string
	Region    string
	SecretID  string
	NameSpace string
	TLSConfig *api.TLSConfig
}

func DefaultClientConfig() ClientConfig {
	clientConf := ClientConfig{
		Address:   os.Getenv("NOMAD_ADDR"),
		Region:    "",
		SecretID:  "",
		NameSpace: "",
		TLSConfig: nil,
	}

	return clientConf
}

func NewClient(config *ClientConfig) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address:   config.Address,
		Region:    config.Region,
		SecretID:  config.SecretID,
		Namespace: config.NameSpace,
		TLSConfig: config.TLSConfig,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewStream(config *ClientConfig) *Stream {
	client, _ := NewClient(config)
	return &Stream{
		nomad: client,
		L:     hclog.Default(),
	}
}

func (s *Stream) Subscribe(ctx context.Context) (<-chan *api.Events, error) {
	events := s.nomad.EventStream()

	topics := map[api.Topic][]string{
		api.Topic("Allocation"): {"*"},
	}

	eventCh, err := events.Stream(ctx, topics, 0, &api.QueryOptions{})
	if err != nil {
		s.L.Error("error creating event stream client", "error", err)
		return nil, err
	}

	// Create a channel to return events to the caller
	eventStream := make(chan *api.Events)

	go func() {
		defer close(eventStream)

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventCh:
				if event.Err != nil {
					s.L.Warn("error from event stream", "error", err)
					return
				}
				if event.IsHeartbeat() {
					continue
				}

				// Send the event to the caller
				eventStream <- event
			}
		}
	}()

	return eventStream, nil
}
