package handlers

import (
	"fmt"

	"github.com/devops-rob/boundary-census/clients/boundary"

	"github.com/hashicorp/go-hclog"
)

// ServiceInstance is an abstraction of a Nomad allocation or Kubernetes pod
type ServiceInstance struct {
	// the URI or IP address for the instance
	Location string

	// The Ports exposed by the service
	Ports []uint32
}

// Target is a handler that can create / update / delete Boundary Targets
type Target struct {
	Log            hclog.Logger
	BoundaryClient boundary.Client
}

func NewTarget(l hclog.Logger, b boundary.Client) *Target {
	return &Target{Log: l, BoundaryClient: b}
}

// Create new targets from the given ServiceInstance
func (t *Target) Create(s *ServiceInstance, name, scope, project string) ([]string, error) {
	// attempt to find the project
	project_id, err := t.BoundaryClient.FindProjectIDByName(scope, project)
	if err != nil {
		t.Log.Error("unable to find project", "scope", scope, "project", project, "error", err)
		return nil, fmt.Errorf("unable to find project")
	}

	// create a target for every port
	for _, p := range s.Ports {
		targetName := fmt.Sprintf("%s_%d", name, p)
		_, err = t.BoundaryClient.CreateTarget(targetName, p, project_id)
		if err != nil {
			t.Log.Info("unable to create target", "scope", scope, "project", project, "error", err)
			return nil, fmt.Errorf("unable to create target: %s", err)
		}
	}

	return []string{}, nil
}
