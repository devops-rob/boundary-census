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
