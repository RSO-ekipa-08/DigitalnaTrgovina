terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0"
    }
  }
  required_version = ">= 1.0"
}

provider "kubernetes" {
  host                   = data.terraform_remote_state.aks.outputs.kube_config.0.host
  client_certificate     = base64decode(data.terraform_remote_state.aks.outputs.kube_config.0.client_certificate)
  client_key             = base64decode(data.terraform_remote_state.aks.outputs.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(data.terraform_remote_state.aks.outputs.kube_config.0.cluster_ca_certificate)
} 