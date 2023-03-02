package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, config string) string {
	dir := t.TempDir()
	configfile := filepath.Join(dir, "config.hcl")

	err := os.WriteFile(configfile, []byte(config), os.ModePerm)
	require.NoError(t, err)

	return configfile
}

func TestParsesConfig(t *testing.T) {
	cfg := setup(t, mockConfig)

	c, err := Parse(cfg)
	require.NoError(t, err)
	require.NotNil(t, c)

	require.Equal(t, "http://localhost:4646", c.Nomad.Address)
	require.Equal(t, "abc123", c.Nomad.Token)
	require.Equal(t, "myregion", c.Nomad.Region)
	require.Equal(t, "mynamespace", c.Nomad.Namespace)

	require.Equal(t, "nic", c.Boundary.Username)
	require.Equal(t, "password", c.Boundary.Password)
	require.Equal(t, "http://myaddress.com", c.Boundary.Address)
	require.Equal(t, "myorg", c.Boundary.OrgID)
	require.Equal(t, "123", c.Boundary.AuthMethodID)
	require.Equal(t, "hashicorp", c.Boundary.DefaultProject)
	require.Equal(t, []string{"developers"}, c.Boundary.DefaultGroups)
}

var mockConfig = `
config "controller" {
  nomad {
    address = "http://localhost:4646"
    token = "abc123" 
    region = "myregion"
    namespace = "mynamespace"
  }

  boundary {
    username = "nic"
    password = "password"
    address = "http://myaddress.com"

    org_id = "myorg"
    auth_method_id = "123"
    default_project = "hashicorp"
    default_groups = ["developers"]
  }
}
`
