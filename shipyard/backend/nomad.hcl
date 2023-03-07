variable "cn_network" {
  default = "backend"
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

template "nomad_config_override" {
  source = <<-EOF
  client {
    #{{ if ne .Vars.name "" }}
    host_volume "#{{ .Vars.name }}" {
      path = "#{{ .Vars.path }}"
    }
    #{{ end }}
  }

  vault {
    enabled = true
    address = "http://vault.container.shipyard.run:8200"
    token = "root"
  }
  EOF

  vars = {
    name = var.cn_nomad_client_host_volume.name
    path = var.cn_nomad_client_host_volume.destination
  }

  destination = "${data("nomad_config")}/client_overide.hcl"
}

variable "cn_nomad_client_config" {
  default = "${data("nomad_config")}/client_overide.hcl"
}

module "consul_nomad" {
  source = "github.com/shipyard-run/blueprints?ref=3068f20075208fd43d4b1ba17bbd584e29ba5f2a/modules//consul-nomad"
}

output "NOMAD_ADDR" {
  value = cluster_api("nomad_cluster.local")
}