config "controller" {
  nomad {
    address = env("NOMAD_ADDR")
  }

  boundary {
    username = env("BOUNDARY_USER")
    password = env("BOUNDARY_PASSWORD")
    address = env("BOUNDARY_ADDR")

    org_id = file("./shipyard/files/org_id")
    auth_method_id = file("./shipyard/files/auth_method_id")
    default_project = "myproject"
    default_groups = ["developers"]
  }
}
