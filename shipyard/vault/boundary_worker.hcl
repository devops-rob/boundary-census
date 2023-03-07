container "boundary_worker_vault" {
  image {
    name = "nicholasjackson/boundary-worker-hcp:v0.12.0"
  }

  command = ["/boundary/startup.sh"]

  volume {
    source      = "./files/boundary_worker"
    destination = "/boundary"
  }

  env {
    key   = "worker_name"
    value = "vault"
  }

  env {
    key   = "cluster_id"
    value = var.boundary_cluster_id
  }

  env {
    key   = "username"
    value = var.boundary_username
  }

  env {
    key   = "password"
    value = var.boundary_password
  }

  env {
    key   = "auth_method_id"
    value = var.boundary_auth_method_id
  }
}

variable "boundary_cluster_id" {
  default = ""
}

variable "boundary_username" {
  default = ""
}

variable "boundary_password" {
  default = ""
}

variable "boundary_auth_method_id" {
  default = ""
}