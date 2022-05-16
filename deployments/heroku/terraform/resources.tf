resource "heroku_app" "default" {
  name   = var.heroku_app_name
  region = var.heroku_region
  stack = "container"
}

resource "heroku_addon" "database" {
  app  = heroku_app.default.name
  plan = "heroku-postgresql:hobby-dev"
}

resource "heroku_build" "default" {
  app = heroku_app.default.id
  source {
    path = "../app"
  }
}
