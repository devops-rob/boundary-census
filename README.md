# Boundary Census

## Configuration

The configuration file for the server is specified as HCL, below is an example config file that contains all the possible
fields.

```hcl
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
```

## Setup Local Nomad, Boundary, Consul

To setup and configure a local Nomad, Boundary, and Consul server use the following command:

```shell
shipyard run ./shipyard
```

You can determine the local addresses for the Boundary, Consul, and Nomad clusters by running:

```
shipyard output
```

These can also be set as environment variables with the following command:

```
eval $(shipyard env)
```

## Running

To run the server you can use the following command:

```shell
go run main.go -config=./example_config.hcl
```
