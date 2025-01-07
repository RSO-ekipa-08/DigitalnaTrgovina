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
    AUTH0_DOMAIN        = var.auth_AUTH0_DOMAIN
    AUTH0_CALLBACK_URL  = var.auth_AUTH0_CALLBACK_URL
  }
}

# Secret
resource "kubernetes_secret" "rso_secrets" {
  metadata {
    name      = "rso-secrets"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  data = {
    DB_PASSWORD         = var.postgresql_password
    POSTGRES_PASSWORD   = var.postgresql_password
    JWT_SECRET          = var.jwt_secret
    AUTH0_CLIENT_ID     = var.auth_AUTH0_CLIENT_ID
    AUTH0_CLIENT_SECRET = var.auth_AUTH0_CLIENT_SECRET
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
            container_port = 3000
          }

          port {
            container_port = 50051
          }

          env {
            name = "AUTH0_CLIENT_ID"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.rso_secrets.metadata[0].name
                key  = "AUTH0_CLIENT_ID"
              }
            }
          }

          env {
            name  = "AUTH0_DOMAIN"
            value = var.auth_AUTH0_DOMAIN
          }

          env {
            name = "AUTH0_CLIENT_SECRET"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.rso_secrets.metadata[0].name
                key  = "AUTH0_CLIENT_SECRET"
              }
            }
          }

          env {
            name  = "AUTH0_CALLBACK_URL"
            value = var.auth_AUTH0_CALLBACK_URL
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
      name        = "grpc"
      port        = 50051
      target_port = 50051
    }

    port {
      name        = "http"
      port        = 80
      target_port = 3000
    }

    type = "LoadBalancer"
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

          port {
            container_port = 50051
          }

          env {
            name  = "POSTGRES_HOST"
            value = data.terraform_remote_state.aks.outputs.postgresql_server_fqdn
          }

          env {
            name  = "DATABASE_URL"
            value = "postgres://psqladmin:${var.postgresql_password}@${data.terraform_remote_state.aks.outputs.postgresql_server_fqdn}:5432/reviews_db"
          }

          env {
            name  = "POSTGRES_PORT"
            value = "5432"
          }

          env {
            name  = "POSTGRES_USER"
            value = "psqladmin"
          }

          env {
            name  = "POSTGRES_DB"
            value = "reviews_db"
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

        # Add init container to wait for database
        init_container {
          name  = "wait-for-db"
          image = "postgres:16"

          command = ["/bin/sh", "-c"]
          args = [
            "until pg_isready -h ${data.terraform_remote_state.aks.outputs.postgresql_server_fqdn} -p 5432 -U psqladmin; do echo waiting for database; sleep 2; done;"
          ]

          env {
            name = "PGPASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.rso_secrets.metadata[0].name
                key  = "POSTGRES_PASSWORD"
              }
            }
          }
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
