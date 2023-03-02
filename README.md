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

## Running

To run the server you can use the following command:

```shell
go run main.go -config=./myconfig.hcl
```
