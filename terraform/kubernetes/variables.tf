variable "jwt_secret" {
  description = "Secret za JWT tokene"
  type        = string
  sensitive   = true
}

# variable "postgresql_password" {
#   description = "Geslo za PostgreSQL administratorja"
#   type        = string
#   sensitive   = true
# }

variable "reviews_DATABASE_URL" {
  description = "URL za PostgreSQL bazo za reviews."
  type        = string
  sensitive   = true
}

variable "auth_AUTH0_CLIENT_ID" {
  description = "Auth0 client ID"
  type        = string
  sensitive   = true
}

variable "auth_AUTH0_DOMAIN" {
  description = "Auth0 domain"
  type        = string
  sensitive   = true
}

variable "auth_AUTH0_CLIENT_SECRET" {
  description = "Auth0 client secret"
  type        = string
  sensitive   = true
}

variable "auth_AUTH0_CALLBACK_URL" {
  description = "Auth0 callback URL"
  type        = string
  sensitive   = true
}

variable "app_DATABASE_URL" {
  description = "URL za podatkovno bazo app-service"
  type        = string
  sensitive   = true
}

variable "app_STORAGE_ENDPOINT" {
  description = "Končna točka za object storage"
  type        = string
  sensitive   = true
}

variable "app_STORAGE_ACCESS_KEY" {
  description = "Dostopni ključ za object storage"
  type        = string
  sensitive   = true
  default     = ""
}

variable "app_STORAGE_SECRET_KEY" {
  description = "Skrivni ključ za object storage"
  type        = string
  sensitive   = true
}
