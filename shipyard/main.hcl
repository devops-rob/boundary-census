module "frontend" {
  source = "./frontend"
}

module "backend" {
  source = "./backend"
}

module "vault" {
  source = "./vault"
}