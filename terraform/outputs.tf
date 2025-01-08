output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

# output "postgresql_server_name" {
#   value = azurerm_postgresql_flexible_server.postgresql.name
# }

# output "postgresql_server_fqdn" {
#   value = azurerm_postgresql_flexible_server.postgresql.fqdn
# }

output "kube_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_config
  sensitive = true
}

output "grafana_endpoint" {
  description = "URL naslov Grafane"
  value       = azurerm_dashboard_grafana.grafana.endpoint
}

output "log_analytics_workspace_id" {
  description = "ID Log Analytics workspace-a"
  value       = azurerm_log_analytics_workspace.aks.id
} 