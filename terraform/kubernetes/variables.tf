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