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

variable "stripe_secret_key" {
  description = "Stripe Secret Key"
  type        = string
  sensitive   = true
}

variable "payment_success_url" {
  description = "URL for successful payments"
  type        = string
}

variable "payment_cancel_url" {
  description = "URL for cancelled payments"
  type        = string
}

variable "app_tenants" {
  description = "Konfiguracija za različne najemnike app-service"
  type = map(object({
    database_url = string
    storage_endpoint = string
    storage_access_key = optional(string, "")
    storage_secret_key = string
    replicas = optional(number, 1)
    cpu_request = optional(string, "100m")
    memory_request = optional(string, "128Mi")
    cpu_limit = optional(string, "200m")
    memory_limit = optional(string, "256Mi")
  }))
  sensitive = true
  validation {
    # At least one tenant must be defined
    condition     = length(var.app_tenants) > 0
    error_message = "At least one tenant must be defined."
  }
}
