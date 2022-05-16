provider "heroku" {
  email   = var.heroku_email
  api_key = var.heroku_api_key
}

output "web_url" {
  value = heroku_app.default.web_url 
}
