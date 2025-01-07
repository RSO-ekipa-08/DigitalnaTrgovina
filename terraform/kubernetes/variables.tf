variable "jwt_secret" {
  description = "Secret za JWT tokene"
  type        = string
  sensitive   = true
}

variable "postgresql_password" {
  description = "Geslo za PostgreSQL administratorja"
  type        = string
  sensitive   = true
}

variable "reviews_DATABASE_URL" {
  description = "URL za PostgreSQL bazo za reviews."
  type        = string
  sensitive   = true
}

variable "reviews_POSTGRES_HOST" {
  description = "Host za PostgreSQL bazo za reviews."
  type        = string
  sensitive   = true
}

variable "reviews_POSTGRES_PORT" {
  description = "Port za PostgreSQL bazo za reviews."
  type        = string
  sensitive   = true
}

variable "reviews_POSTGRES_USER" {
  description = "Uporabniško ime za PostgreSQL bazo za reviews."
  type        = string
  sensitive   = true
}

variable "reviews_POSTGRES_PASSWORD" {
  description = "Geslo za PostgreSQL bazo za reviews."
  type        = string
  sensitive   = true
}

variable "reviews_POSTGRES_DB" {
  description = "Naziv baze za PostgreSQL bazo za reviews."
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
