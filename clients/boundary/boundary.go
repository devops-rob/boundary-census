package boundary

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/authmethods"
	"github.com/hashicorp/boundary/api/hostcatalogs"
	"github.com/hashicorp/boundary/api/hosts"
	"github.com/hashicorp/boundary/api/hostsets"
	"github.com/hashicorp/boundary/api/scopes"
	"github.com/hashicorp/boundary/api/targets"
	//"github.com/hashicorp/nomad/helper/authmethods"
)

//go:generate mockery --name Client
type Client interface {
	// CreateTarget creates a new Boundary target with the given options
	CreateTarget(name string, address string, port uint32, scopeId, ingressFilter, egressFilter string) (*targets.Target, error)

	// FindProjectIDByName attempts to find a project in an organization
	// project can be referenced by either the name or the id of the project
	// if a project is found the id and a nil error is returned
	// if a project is not found a ProjectNotFound error is returned
	FindProjectIDByName(org, name string) (string, error)

	// DeleteTargetsWithPrefix deletes all targets in the scope that have the given prefix
	DeleteTargetsWithPrefix(prefix, scopeId string) error
}

// ProjectNotFoundError is returned by FindProjectIDByName when the given project
// is not found in the organization
var (
	ProjectNotFoundError = fmt.Errorf("project not found")
	TargetNotFoundError  = fmt.Errorf("target not found")
)

type ClientImpl struct {
	*api.Client
	organization string
	scope        string
	isEnterprise bool
}

func New(address string, organization string, scope string, authmethod string, credentials map[string]interface{}, isEnterprise bool) (Client, error) {
	config := &api.Config{
		Addr: address,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create client from config: %s", err)
	}

	amClient := authmethods.NewClient(client)
	authenticationResult, err := amClient.Authenticate(context.Background(), authmethod, "login", credentials)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate client: %s", err)
	}

	client.SetToken(fmt.Sprint(authenticationResult.Attributes["token"]))

	return &ClientImpl{client, organization, scope, isEnterprise}, nil
}

func (c *ClientImpl) FindProjectIDByName(org, name string) (string, error) {
	var opts []scopes.Option
	opts = append(opts, scopes.WithRecursive(true))

	client := scopes.NewClient(c.Client)

	result, err := client.List(context.Background(), org, opts...)
	if err != nil {
		return "", err
	}

	for _, scope := range result.Items {
		if scope.Id == name && scope.Type == "project" {
			return scope.Id, nil
		}

		if scope.Name == name && scope.Type == "project" {
			return scope.Id, nil
		}
	}

	return "", ProjectNotFoundError
}

func (c *ClientImpl) GetHostCatalogByName(name string, scopeId string) (*hostcatalogs.HostCatalog, error) {
	var opts []hostcatalogs.Option

	client := hostcatalogs.NewClient(c.Client)

	result, err := client.List(context.Background(), scopeId, opts...)
	if err != nil {
		return nil, err
	}

	for _, hostCatalog := range result.Items {
		if hostCatalog.Name == name {
			return hostCatalog, nil
		}
	}

	return nil, fmt.Errorf("host catalog not found")
}

func (c *ClientImpl) CreateHostCatalog(name string, scopeId string) (*hostcatalogs.HostCatalog, error) {
	var opts []hostcatalogs.Option

	opts = append(opts, hostcatalogs.WithName(name))

	client := hostcatalogs.NewClient(c.Client)

	result, err := client.Create(context.Background(), "static", scopeId, opts...)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) DeleteHostCatalog(id string, scopeId string) error {
	var opts []hostcatalogs.Option

	client := hostcatalogs.NewClient(c.Client)

	_, err := client.Delete(context.Background(), id, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) GetHostsetByName(name string, hostCatalogId string) (*hostsets.HostSet, error) {
	var opts []hostsets.Option

	client := hostsets.NewClient(c.Client)

	result, err := client.List(context.Background(), hostCatalogId, opts...)
	if err != nil {
		return nil, err
	}

	for _, hostSet := range result.Items {
		if hostSet.Name == name {
			return hostSet, nil
		}
	}

	return nil, fmt.Errorf("host set not found")
}

func (c *ClientImpl) CreateHostset(name string, hostCatalogId string) (*hostsets.HostSet, error) {
	var opts []hostsets.Option

	opts = append(opts, hostsets.WithName(name))

	client := hostsets.NewClient(c.Client)

	result, err := client.Create(context.Background(), hostCatalogId, opts...)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) DeleteHostset(hostSetId string) error {
	var opts []hostsets.Option

	client := hostsets.NewClient(c.Client)

	_, err := client.Delete(context.Background(), hostSetId, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) GetHostByName(name string, hostCatalogId string) (*hosts.Host, error) {
	var opts []hosts.Option

	client := hosts.NewClient(c.Client)

	result, err := client.List(context.Background(), hostCatalogId, opts...)
	if err != nil {
		return nil, err
	}

	for _, host := range result.Items {
		if host.Name == name {
			return host, nil
		}
	}

	return nil, fmt.Errorf("host not found")
}

func (c *ClientImpl) CreateHost(name string, address string, hostCatalogId string) (*hosts.Host, error) {
	var opts []hosts.Option

	opts = append(opts, hosts.WithName(name))
	opts = append(opts, hosts.WithStaticHostAddress(address))

	client := hosts.NewClient(c.Client)

	result, err := client.Create(context.Background(), hostCatalogId, opts...)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) UpdateHost(name string, address string, id string, version uint32) (*hosts.Host, error) {
	var opts []hosts.Option

	opts = append(opts, hosts.WithName(name))
	opts = append(opts, hosts.WithStaticHostAddress(address))

	client := hosts.NewClient(c.Client)

	result, err := client.Update(context.Background(), id, version, opts...)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) DeleteHost(id string) error {
	var opts []hosts.Option

	client := hosts.NewClient(c.Client)

	_, err := client.Delete(context.Background(), id, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) GetTargetByName(name string, scopeId string) (*targets.Target, error) {
	var opts []targets.Option

	client := targets.NewClient(c.Client)

	result, err := client.List(context.Background(), scopeId, opts...)
	if err != nil {
		return nil, err
	}

	for _, target := range result.Items {
		if target.Name == name {
			return target, nil
		}
	}

	return nil, TargetNotFoundError
}

func (c *ClientImpl) CreateTarget(name string, address string, port uint32, scopeId, ingressFilter, egressFilter string) (*targets.Target, error) {
	// first check if the target exists
	t, err := c.GetTargetByName(name, scopeId)
	if err != nil && err != TargetNotFoundError {
		return nil, err
	}

	client := targets.NewClient(c.Client)

	var opts []targets.Option
	opts = append(opts, targets.WithTcpTargetDefaultPort(port))
	opts = append(opts, targets.WithAddress(address))
	opts = append(opts, targets.WithName(name))

	if c.isEnterprise {
		if ingressFilter != "" {
			opts = append(opts, targets.WithIngressWorkerFilter(ingressFilter))
		}

		if egressFilter != "" {
			opts = append(opts, targets.WithEgressWorkerFilter(egressFilter))
		}
	} else {
		// if oss
		if egressFilter != "" {
			opts = append(opts, targets.WithEgressWorkerFilter(egressFilter))
		}
	}

	// create
	if err == TargetNotFoundError {
		result, err := client.Create(context.Background(), "tcp", scopeId, opts...) // check resource type
		if err != nil {
			return nil, err
		}

		return result.Item, nil
	}

	opts = append(opts, targets.WithAutomaticVersioning(true))

	// update
	result, err := client.Update(context.Background(), t.Id, 0, opts...)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) CreateTargetWithHost(name string, port uint32, scopeId string, hostId string) (*targets.Target, error) {
	var opts []targets.Option
	opts = append(opts, targets.WithTcpTargetDefaultPort(port))
	opts = append(opts, targets.WithName(name))
	opts = append(opts, targets.WithHostId(hostId))

	client := targets.NewClient(c.Client)

	result, err := client.Create(context.Background(), "tcp", scopeId, opts...) // check resource type
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) UpdateTarget(name string, port uint32, id string, hostId string, version uint32) (*targets.Target, error) {
	var opts []targets.Option
	opts = append(opts, targets.WithTcpTargetDefaultPort(port))
	opts = append(opts, targets.WithName(name))
	opts = append(opts, targets.WithHostId(hostId))

	client := targets.NewClient(c.Client)

	result, err := client.Update(context.Background(), id, version, opts...) // check resource type
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (c *ClientImpl) DeleteTargetsWithPrefix(prefix, scopeId string) error {
	var opts []targets.Option
	client := targets.NewClient(c.Client)

	result, err := client.List(context.Background(), scopeId, opts...)
	if err != nil {
		return err
	}

	for _, target := range result.Items {
		if strings.HasPrefix(target.Name, prefix) {
			_, err := client.Delete(context.Background(), target.Id, opts...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ClientImpl) DeleteTarget(id string) error {
	var opts []targets.Option

	client := targets.NewClient(c.Client)

	_, err := client.Delete(context.Background(), id, opts...)
	if err != nil {
		return err
	}

	return nil
}
