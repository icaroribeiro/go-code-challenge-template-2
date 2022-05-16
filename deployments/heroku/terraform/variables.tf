variable "heroku_email" {
  description = "Heroku email account"
}

variable "heroku_api_key" {
  description = "Heroku authorization token"
}

variable "heroku_app_name" {
  description = "Heroku application name"
}

variable "heroku_region" {
  description = "Heroku application region"
  default = "us"
}
