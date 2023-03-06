network "local" {
    subnet = "10.0.0.0/16"
}

container "boundary_db" {
  network {
      name = "network.local"
  }

  image {
      name = "postgres:15.1"
  }

  env {
    key = "POSTGRES_USER"
    value = "postgres"
  }

  env {
    key = "POSTGRES_PASSWORD"
    value = "postgres"
  }

  env {
    key = "POSTGRES_DB"
    value ="boundary"
  }

  port {
    host   = 5432
    local  = 5432
    remote = 5432
  }
}

template "boundary_db_init" {
  source = <<-EOF
    #! /bin/sh

    command="boundary database init -config /boundary/config/boundary_server.hcl -skip-auth-method-creation -skip-host-resources-creation -skip-scopes-creation -skip-target-creation"

    # Wait until the database can be contacted
    while ! $command
    do
      echo "waiting for db server"
      sleep 10
    done
  EOF

  destination = "${data("temp")}/boundary_init.sh"
}

exec_remote "boundary_db_init" {
  depends_on = ["container.boundary_db", "template.boundary_db_init"]

  image {
      name = "hashicorp/boundary:0.12"
  }

  network {
    name = "network.local"
  }

  cmd = "sh"
  args = [
    "/files/boundary_init.sh"
  ]

  volume {
    source = "./files/config"
    destination = "/boundary/config"
  }

  volume {
    source = data("temp")
    destination = "/files"
  }
}

container "boundary" {
  depends_on = ["exec_remote.boundary_db_init"]

  network {
    name = "network.local"
    ip_address = "10.0.0.200"
  }

  image {
      name = "hashicorp/boundary:0.12"
  }

  command = ["boundary", "server", "-config", "/boundary/config/boundary_server.hcl"]

  volume {
    source = "./files/config"
    destination = "/boundary/config"
  }

  port {
    host   = 9200
    local  = 9200
    remote = 9200
  }

  port {
    host   = 9201
    local  = 9201
    remote = 9201
  }

  port {
    host   = 9202
    local  = 9202
    remote = 9202
  }

  port {
    host   = 9203
    local  = 9203
    remote = 9203
  }
}

exec_remote "boundary_init" {
  depends_on = ["container.boundary"]

  image   {
    name = "shipyardrun/hashicorp-tools:v0.11.0"
  }

  network {
    name = "network.local"
  }

  cmd = "/bin/bash"
  args = ["/files/setup_boundary.sh"]

  volume {
    source      = "./files"
    destination = "/files"
  }

  env {
    key = "BOUNDARY_ADDR"
    value = "http://boundary.container.shipyard.run:9200"
  }

  working_directory = "/files"
}

output "BOUNDARY_ADDR" {
  value = "http://boundary.container.shipyard.run:9200"
}

output "BOUNDARY_USER" {
  value = "nicj"
}

output "BOUNDARY_PASSWORD" {
  value = "password"
}
