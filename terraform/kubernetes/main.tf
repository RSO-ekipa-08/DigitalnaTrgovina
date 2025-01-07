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
    DB_HOST             = ""
    DB_PORT             = "5432"
    DB_USER             = "psqladmin"
    DB_NAME             = "rso"
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
    DB_PASSWORD = var.postgresql_password
    JWT_SECRET  = var.jwt_secret
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

# Info Server
resource "kubernetes_deployment" "info" {
  metadata {
    name      = "info-server"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "info-server"
      }
    }

    template {
      metadata {
        labels = {
          app = "info-server"
        }
      }

      spec {
        container {
          name  = "info-server"
          image = "nginxdemos/hello:latest"

          port {
            container_port = 80
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
        }
      }
    }
  }
}

resource "kubernetes_service" "info" {
  metadata {
    name      = "info-server"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.info.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 80  # Zunanji port nastavimo na 80 za HTTP
      target_port = 80  # Notranji port kontejnerja
    }

    type = "LoadBalancer"
  }
}
