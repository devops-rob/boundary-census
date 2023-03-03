variable "cn_network" {
  default = "local"
}

variable "cn_nomad_cluster_name" {
  default = "nomad.local"
}

variable "cn_nomad_version" {
  default = "1.4.0"
}

variable "cn_nomad_version" {
  default = "1.4.0"
}

variable "cd_consul_version" {
  default = "1.15.0"
}

module "consul_nomad" {
  source = "github.com/shipyard-run/blueprints?ref=3068f20075208fd43d4b1ba17bbd584e29ba5f2a/modules//consul-nomad"
}

output "NOMAD_ADDR" {
  value = cluster_api("nomad_cluster.local")
}
