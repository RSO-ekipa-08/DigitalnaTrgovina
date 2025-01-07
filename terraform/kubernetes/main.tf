# Namespace
resource "kubernetes_namespace" "rso" {
  metadata {
    name = "rso"
  }
}

# ConfigMap
resource "kubernetes_config_map" "rso_config" {
  metadata {
    name      = "rso-config"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  data = {
    ENVIRONMENT         = "prod"
    DB_HOST             = data.terraform_remote_state.aks.outputs.postgresql_server_fqdn
    DB_PORT             = "5432"
    DB_USER             = "psqladmin"
    DB_NAME             = "rso"
    POSTGRES_USER       = "psqladmin"
    POSTGRES_DB         = "rso"
    AUTH_SERVICE_URL    = "http://auth-service:8080"
    PAYMENT_SERVICE_URL = "http://payment-service:8080"
    REVIEWS_SERVICE_URL = "http://reviews-service:8080"
    APP_SERVICE_URL     = "http://app-service:8080"
  }
}

# Secret
resource "kubernetes_secret" "rso_secrets" {
  metadata {
    name      = "rso-secrets"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  data = {
    DB_PASSWORD       = var.postgresql_password
    POSTGRES_PASSWORD = var.postgresql_password
    JWT_SECRET        = var.jwt_secret
  }
}

# Auth Service
resource "kubernetes_deployment" "auth" {
  metadata {
    name      = "auth-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "auth-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "auth-service"
        }
      }

      spec {
        container {
          name  = "auth-service"
          image = "ghcr.io/rso-ekipa-08/authentication:latest"

          port {
            container_port = 8080
          }

          env_from {
            config_map_ref {
              name = kubernetes_config_map.rso_config.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.rso_secrets.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "200m"
              memory = "256Mi"
            }
          }

          #   readiness_probe {
          #     http_get {
          #       path = "/health"
          #       port = 8080
          #     }
          #     initial_delay_seconds = 5
          #     period_seconds       = 10
          #   }
        }
      }
    }
  }
}

resource "kubernetes_service" "auth" {
  metadata {
    name      = "auth-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.auth.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 8080
      target_port = 8080
    }

    type = "ClusterIP"
  }
}

# Payment Service
resource "kubernetes_deployment" "payment" {
  metadata {
    name      = "payment-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "payment-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "payment-service"
        }
      }

      spec {
        container {
          name  = "payment-service"
          image = "ghcr.io/rso-ekipa-08/payment:latest"

          port {
            container_port = 8080
          }

          env_from {
            config_map_ref {
              name = kubernetes_config_map.rso_config.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.rso_secrets.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "200m"
              memory = "256Mi"
            }
          }

          #   readiness_probe {
          #     http_get {
          #       path = "/health"
          #       port = 8080
          #     }
          #     initial_delay_seconds = 5
          #     period_seconds       = 10
          #   }
        }
      }
    }
  }
}

resource "kubernetes_service" "payment" {
  metadata {
    name      = "payment-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.payment.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 8080
      target_port = 8080
    }

    type = "ClusterIP"
  }
}

# Reviews Service
resource "kubernetes_deployment" "reviews" {
  metadata {
    name      = "reviews-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "reviews-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "reviews-service"
        }
      }

      spec {
        container {
          name  = "reviews-service"
          image = "ghcr.io/rso-ekipa-08/reviews:latest"

          env_from {
            config_map_ref {
              name = kubernetes_config_map.rso_config.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.rso_secrets.metadata[0].name
            }
          }

          env {
            name  = "POSTGRES_HOST"
            value = data.terraform_remote_state.aks.outputs.postgresql_server_fqdn
          }

          env {
            name  = "POSTGRES_USER"
            value = "psqladmin"
          }

          env {
            name  = "POSTGRES_DB"
            value = "reviews_db" # Match the database name we created
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.rso_secrets.metadata[0].name
                key  = "POSTGRES_PASSWORD"
              }
            }
          }

          port {
            container_port = 50051 # This matches your Rust service port
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "200m"
              memory = "256Mi"
            }
          }

          #   readiness_probe {
          #     http_get {
          #       path = "/health"
          #       port = 8080
          #     }
          #     initial_delay_seconds = 5
          #     period_seconds       = 10
          #   }
        }
      }
    }
  }
}

resource "kubernetes_service" "reviews" {
  metadata {
    name      = "reviews-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.reviews.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 50051
      target_port = 50051
    }

    type = "ClusterIP"
  }
}

# App Service
resource "kubernetes_deployment" "app" {
  metadata {
    name      = "app-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "app-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "app-service"
        }
      }

      spec {
        container {
          name  = "app-service"
          image = "ghcr.io/rso-ekipa-08/app-service:latest"

          port {
            container_port = 8080
          }

          env_from {
            config_map_ref {
              name = kubernetes_config_map.rso_config.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.rso_secrets.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "200m"
              memory = "256Mi"
            }
          }

          #   readiness_probe {
          #     http_get {
          #       path = "/health"
          #       port = 8080
          #     }
          #     initial_delay_seconds = 5
          #     period_seconds       = 10
          #   }
        }
      }
    }
  }
}

resource "kubernetes_service" "app" {
  metadata {
    name      = "app-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.app.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 8080
      target_port = 8080
    }

    type = "LoadBalancer"
  }
}
