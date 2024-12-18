output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "kubernetes_cluster_name" {
  value = azurerm_kubernetes_cluster.aks.name
}

output "postgresql_server_name" {
  value = azurerm_postgresql_flexible_server.postgresql.name
}

output "postgresql_server_fqdn" {
  value = azurerm_postgresql_flexible_server.postgresql.fqdn
}

# Kubeconfig za dostop do AKS
output "aks_credentials" {
  value     = azurerm_kubernetes_cluster.aks.kube_config_raw
  sensitive = true
} 