exec_remote "get_plugin" {
  image {
    name = "shipyardrun/tools:v0.7.0"
  }

  volume {
    source      = "./files/"
    destination = "/files"
  }

  cmd = "/files/get-plugins.sh"
}

container "vault" {
  network {
    name       = "network.frontend"
    ip_address = "10.0.1.210"
  }

  network {
    name       = "network.backend"
    ip_address = "10.0.2.210"
  }

  network {
    name       = "network.vault"
    ip_address = "10.0.3.210"
  }

  image {
    name = "hashicorp/vault:1.12.3"
  }

  port {
    host   = 8200
    local  = 8200
    remote = 8200
  }

  env {
    key   = "VAULT_DEV_ROOT_TOKEN_ID"
    value = "root"
  }

  env {
    key   = "VAULT_ADDR"
    value = "http://localhost:8200"
  }

  env {
    key   = "VAULT_TOKEN"
    value = "root"
  }

  privileged = true

  command = [
    "vault",
    "server",
    "-dev",
    "-dev-root-token-id=root",
    "-dev-listen-address=0.0.0.0:8200",
    "-dev-plugin-dir=/plugins"
  ]

  volume {
    source      = "./files/vault_plugins"
    destination = "/plugins"
  }

  volume {
    source      = "./files/vault"
    destination = "/files"
  }


  depends_on = ["exec_remote.get_plugin"]
}

output "VAULT_ADDR" {
  value = "http://localhost:8200"
}

output "VAULT_TOKEN" {
  value = "root"
}