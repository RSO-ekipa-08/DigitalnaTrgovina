variable "project_name" {
  description = "Ime projekta, ki se uporablja kot prefix za vse resurse"
  type        = string
  default     = "rso"
}

variable "environment" {
  description = "Okolje (dev, prod)"
  type        = string
  default     = "prod"
}

variable "location" {
  description = "Azure regija za vse resurse"
  type        = string
  default     = "germanywestcentral"
}

variable "kubernetes_version" {
  description = "Verzija Kubernetes-a"
  type        = string
  default     = "1.31"
}

variable "node_count" {
  description = "Število node-ov v AKS cluster-ju"
  type        = number
  default     = 1
}

variable "node_size" {
  description = "Velikost node-ov v AKS cluster-ju"
  type        = string
  default     = "Standard_B2ms"
}

variable "postgresql_sku" {
  description = "SKU za PostgreSQL"
  type        = string
  default     = "B_Standard_B1ms"
}

variable "postgresql_storage" {
  description = "Velikost storage-a za PostgreSQL v MB"
  type        = number
  default     = 32768  # 32GB
}

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

variable "subscription_id" {
  description = "Azure subscription ID"
  type        = string
  sensitive   = true
}
