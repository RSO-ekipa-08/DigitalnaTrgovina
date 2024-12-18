output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "postgresql_server_name" {
  value = azurerm_postgresql_flexible_server.postgresql.name
}

output "postgresql_server_fqdn" {
  value = azurerm_postgresql_flexible_server.postgresql.fqdn
}

output "kube_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_config
  sensitive = true
} 