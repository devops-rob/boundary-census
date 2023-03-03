disable_mlock = true
log_level  = "debug"
log_format = "standard"

listener "tcp" {
  address = "0.0.0.0:9200"
  purpose = "api"
  tls_disable = true
}

listener "tcp" {
  address = "10.0.0.200:9201"
  purpose = "cluster"
}

listener "tcp" {
  address = "10.0.0.200:9202"
  purpose = "proxy"
}

listener "tcp" {
    purpose = "ops"
    tls_disable = true
}

controller {
  name = "example-controller"
  description = "An example controller"

  database {
    url = "postgresql://postgres:postgres@boundary-db.container.shipyard.run:5432/boundary?sslmode=disable"
    max_open_connections = 5
  }
}

kms "aead" {
    purpose   = "root"
    aead_type = "aes-gcm"
    key       = "sP1fnF5Xz85RrXyELHFeZg9Ad2qt4Z4bgNHVGtD6ung="
    key_id    = "global_root"
}

kms "aead" {
    purpose   = "worker-auth"
    aead_type = "aes-gcm"
    key       = "8fZBjCUfN0TzjEGLQldGY4+iE9AkOvCfjh7+p0GtRBQ="
    key_id    = "global_worker-auth"
}

kms "aead" {
    purpose   = "recovery"
    aead_type = "aes-gcm"
    key       = "8fZBjCUfN0TzjEGLQldGY4+iE9AkOvCfjh7+p0GtRBQ="
    key_id    = "global_recovery"
}

worker {
  name = "server-worker"
  auth_storage_path = "/boundary/worker"

  public_addr = "boundary.container.shipyard.run"

  tags {
    type = ["worker"]
  }
}
