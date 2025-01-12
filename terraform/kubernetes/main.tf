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

          env {
            name  = "APP_SERVICE_URL"
            value = "http://${kubernetes_deployment.app.metadata[0].name}:8080"
          }

          env {
            name  = "REVIEWS_SERVICE_URL"
            value = "http://${kubernetes_deployment.reviews.metadata[0].name}:8080"
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

          liveness_probe {
            http_get {
              path = "/health"
              port = 3000
            }
            initial_delay_seconds = 15
          }
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

# Api Gateway
resource "kubernetes_deployment" "api-gateway" {
  metadata {
    name      = "api-gateway-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "api-gateway-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "api-gateway-service"
        }
      }

      spec {
        container {
          name  = "api-gateway-service"
          image = "ghcr.io/rso-ekipa-08/api-gateway:latest"

          port {
            container_port = 3000
          }

          env {
            name  = "MINIAUTH_ADDRESS"
            value = "auth-service.rso.svc.cluster.local:50051" # Points to the gRPC auth service
          }

          env {
            name  = "APP_SERVICE_URL"
            value = "http://${kubernetes_deployment.app.metadata[0].name}:8080"
          }

          env {
            name  = "REVIEWS_SERVICE_URL"
            value = "http://${kubernetes_deployment.reviews.metadata[0].name}:8080"
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

          liveness_probe {
            http_get {
              path = "/health"
              port = 3000
            }
            initial_delay_seconds = 15
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "api-gateway" {
  metadata {
    name      = "api-gateway-service"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.api-gateway.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 80
      target_port = 3000
    }

    type = "LoadBalancer"
  }
}
# RabbitMQ
resource "kubernetes_deployment" "rabbitmq" {
  metadata {
    name      = "rabbitmq"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "rabbitmq"
      }
    }

    template {
      metadata {
        labels = {
          app = "rabbitmq"
        }
      }

      spec {
        container {
          name  = "rabbitmq"
          image = "rabbitmq:3.12-management"

          port {
            name           = "amqp"
            container_port = 5672
          }

          port {
            name           = "management"
            container_port = 15672
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "256Mi"
            }
            limits = {
              cpu    = "200m"
              memory = "512Mi"
            }
          }

          readiness_probe {
            tcp_socket {
              port = 5672
            }
            initial_delay_seconds = 10
            period_seconds        = 30
          }

          liveness_probe {
            tcp_socket {
              port = 5672
            }
            initial_delay_seconds = 30
            period_seconds        = 30
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "rabbitmq" {
  metadata {
    name      = "rabbitmq"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.rabbitmq.spec[0].template[0].metadata[0].labels.app
    }

    port {
      name        = "amqp"
      port        = 5672
      target_port = 5672
    }

    port {
      name        = "management"
      port        = 15672
      target_port = 15672
    }

    type = "ClusterIP"
  }
}

# Payment Service v2
resource "kubernetes_deployment" "payment_v2" {
  metadata {
    name      = "payment-service-v2"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "payment-service-v2"
      }
    }

    template {
      metadata {
        labels = {
          app = "payment-service-v2"
        }
      }

      spec {
        container {
          name  = "payment-service-v2"
          image = "ghcr.io/rso-ekipa-08/payment_v2:latest"

          port {
            container_port = 3000
          }

          env {
            name = "STRIPE_SECRET_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.payment_v2_secrets.metadata[0].name
                key  = "STRIPE_SECRET_KEY"
              }
            }
          }

          env {
            name = "SUCCESS_URL"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.payment_v2_secrets.metadata[0].name
                key  = "SUCCESS_URL"
              }
            }
          }

          env {
            name = "CANCEL_URL"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.payment_v2_secrets.metadata[0].name
                key  = "CANCEL_URL"
              }
            }
          }

          env {
            name  = "RABBITMQ_URL"
            value = "amqp://rabbitmq:5672"
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

          readiness_probe {
            http_get {
              path = "/health"
              port = 3000
            }
            initial_delay_seconds = 5
            period_seconds        = 10
          }
        }
      }
    }
  }
}

# Payment v2 Secrets
resource "kubernetes_secret" "payment_v2_secrets" {
  metadata {
    name      = "payment-v2-secrets"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  data = {
    STRIPE_SECRET_KEY = var.stripe_secret_key
    SUCCESS_URL       = var.payment_success_url
    CANCEL_URL        = var.payment_cancel_url
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
            name           = "grpc"
            container_port = 50051
          }

          port {
            name           = "graphql"
            container_port = 8080
          }

          env {
            name  = "DATABASE_URL"
            value = var.reviews_DATABASE_URL
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

          readiness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 10
            period_seconds        = 30
            timeout_seconds       = 5
            success_threshold     = 1
            failure_threshold     = 3
          }

          liveness_probe {
            http_get {
              path = "/ready"
              port = 8080
            }
            initial_delay_seconds = 15
            period_seconds        = 30
            timeout_seconds       = 5
            success_threshold     = 1
            failure_threshold     = 3
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
      name        = "grpc"
      port        = 50051
      target_port = 50051
    }
    port {
      name        = "graphql"
      port        = 8080
      target_port = 8080
    }
    type = "ClusterIP"
  }
}

# App Service Multi-tenant Deployment
resource "kubernetes_deployment" "app_multitenancy" {
  for_each = var.app_tenants

  metadata {
    name      = "app-service-${each.key}"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    replicas = each.value.replicas

    selector {
      match_labels = {
        app = "app-service-${each.key}"
      }
    }

    template {
      metadata {
        labels = {
          app = "app-service-${each.key}"
        }
      }

      spec {
        container {
          name  = "app-service-${each.key}"
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

          env {
            name  = "DATABASE_URL"
            value = each.value.database_url
          }

          env {
            name  = "STORAGE_ENDPOINT"
            value = each.value.storage_endpoint
          }

          env {
            name  = "STORAGE_ACCESS_KEY"
            value = each.value.storage_access_key
          }

          env {
            name  = "STORAGE_SECRET_KEY"
            value = each.value.storage_secret_key
          }

          env {
            name  = "RABBITMQ_URL"
            value = "amqp://rabbitmq:5672"
          }

          resources {
            requests = {
              cpu    = each.value.cpu_request
              memory = each.value.memory_request
            }
            limits = {
              cpu    = each.value.cpu_limit
              memory = each.value.memory_limit
            }
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 5
          }
        }
      }
    }
  }
}

# App Service Multi-tenant Service
resource "kubernetes_service" "app_multitenancy" {
  for_each = var.app_tenants

  metadata {
    name      = "app-internal-${each.key}"
    namespace = kubernetes_namespace.rso.metadata[0].name
  }

  spec {
    selector = {
      app = "app-service-${each.key}"
    }

    port {
      port        = 80
      target_port = 8080
    }

    type = "LoadBalancer"
  }
}

# Reviews Service Output
output "reviews_internal_ip" {
  value = kubernetes_service.reviews.status[0].load_balancer[0].ingress[0].ip
}

# Auth Service Output
output "auth_internal_ip" {
  value = kubernetes_service.auth.status[0].load_balancer[0].ingress[0].ip
}

# Multi-tenant App Service Output
output "app_tenant_ips" {
  value = {
    for tenant, service in kubernetes_service.app_multitenancy : 
    tenant => service.status[0].load_balancer[0].ingress[0].ip
  }
}
