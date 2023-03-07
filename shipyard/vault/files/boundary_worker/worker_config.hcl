  disable_mlock = true
  log_level = "debug"

  hcp_boundary_cluster_id = "739d93f9-7f1c-474d-8524-931ab199eaf8"

  listener "tcp" {
    address = "0.0.0.0:9202"
    purpose = "proxy"
  }

  worker {
    auth_storage_path="./auth_data"

    controller_generated_activation_token = "neslat_2KrAuXDV2hE8825gQmcZiKmNa6BRdR3Cgqb1QJXmmm14DwTVnX2acdM4p7dFZzPEsVuiNQPx7iyKmb4c22sr5zniVyXcy"
  
    tags {
      type   = ["raspberypi"]
    }
  }
