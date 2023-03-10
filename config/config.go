package config

import (
	"fmt"
	"strings"

	"github.com/shipyard-run/hclconfig"
	"github.com/shipyard-run/hclconfig/types"
)

// Config defines a struct that holds the config for the controller
type Config struct {
	types.ResourceMetadata `hcl:",remain"`

	Nomad    *Nomad    `hcl:"nomad,block"`
	Boundary *Boundary `hcl:"boundary,block"`
}

// Process is called by hclconfig when it finds a Config resource
// you can do validation and cleanup here
func (c *Config) Process() error {
	c.Boundary.DefaultIngressFilter = strings.TrimSpace(c.Boundary.DefaultIngressFilter)
	c.Boundary.DefaultEgressFilter = strings.TrimSpace(c.Boundary.DefaultEgressFilter)

	if ! c.Boundary.Enterprise {
	if c.Boundary.DefaultIngressFilter!= "" {
		return fmt.Errorf("ingress filters are not supported on oss boundary")
	}
}

	return nil
}

// Nomad is configuration specific to the Nomad scheduler
type Nomad struct {
	Address   string `hcl:"address,optional"`
	Token     string `hcl:"token,optional"`
	Region    string `hcl:"region,optional"`
	Namespace string `hcl:"namespace,optional"`
}

// / Boundary is configuration specific to Boundary
type Boundary struct {
	Enterprise     bool     `hcl:"enterprise,optional"`
	OrgID          string   `hcl:"org_id"`
	DefaultProject string   `hcl:"default_project,optional"`
	DefaultGroups  []string `hcl:"default_groups,optional"`

	AuthMethodID string `hcl:"auth_method_id"`
	Username     string `hcl:"username"`
	Password     string `hcl:"password"`
	Address      string `hcl:"address"`

	DefaultIngressFilter string `hcl:"default_ingress_filter,optional"`
	DefaultEgressFilter  string `hcl:"default_egress_filter,optional"`
}

// Parse the given HCL config file and return the Config
func Parse(config string) (*Config, error) {
	p := hclconfig.NewParser(hclconfig.DefaultOptions())
	p.RegisterType("config", &Config{})
	p.RegisterFunction("trim", func(s string) (string, error) {
		return strings.TrimSpace(s), nil
	})

	c := hclconfig.NewConfig()
	err := p.ParseFile(config, c)
	if err != nil {
		return nil, fmt.Errorf("unable to process file: %s, error: %s", config, err)
	}

	r, err := c.FindResourcesByType("config")
	if err != nil {
		return nil, fmt.Errorf("unable to process file: %s, error: %s", config, err)
	}

	if len(r) != 1 {
		return nil, fmt.Errorf("unable to process file: %s, file does not contain a single config resource", config)
	}

	return r[0].(*Config), nil
}
