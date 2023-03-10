config "controller" {
  nomad {
    address = env("NOMAD_ADDR")
  }

  boundary {
    username = env("BOUNDARY_USER")
    password = env("BOUNDARY_PASSWORD")
    address  = env("BOUNDARY_ADDR")

    org_id          = trim(file("./shipyard/files/org_id"))
    auth_method_id  = trim(file("./shipyard/files/auth_method_id"))
    default_project = "myproject"

    default_egress_filter = <<EOF
      "nomad" in "/tags/environment"
    EOF

    default_groups = ["developers"]
  }
}
