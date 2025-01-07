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
