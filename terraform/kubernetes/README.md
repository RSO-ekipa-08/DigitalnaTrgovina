# Kubernetes Terraform Konfiguracija

Terraform konfiguracija za postavitev in upravljanje Kubernetes resursov.

## Struktura

- `main.tf` - Glavna konfiguracija Kubernetes resursov
- `variables.tf` - Definicije spremenljivk
- `providers.tf` - Konfiguracija Kubernetes providerja
- `remote-state.tf` - Konfiguracija oddaljenega state-a
- `terraform.tfvars.example` - Primer konfiguracijskih vrednosti

## Namestitev

1. Kopirajte `terraform.tfvars.example` v `terraform.tfvars`:
```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Nastavite vrednosti v `terraform.tfvars`

3. Inicializirajte Terraform:
```bash
terraform init
```

## Uporaba

### Pregled in apliciranje sprememb

```bash
# Pregled sprememb
terraform plan

# Apliciranje sprememb
terraform apply
```

## Konfiguracijske vrednosti

Glavne spremenljivke v `terraform.tfvars`:

- `namespace` - Kubernetes namespace za aplikacije
- `replicas` - Število replik za posamezne deploymente
- `image_tag` - Verzija Docker slik
- `domain` - Domena za ingress
- `environment` - Okolje (dev/prod)

## Resursi

Konfiguracija postavlja naslednje Kubernetes resurse:

- Namespace
- Deployments za mikrostoritve
- Services
- Ingress
- ConfigMaps
- Secrets
- Service Accounts

## Odpravljanje težav

Za preverjanje stanja resursov:

```bash
# Izpis stanja
terraform show

# Seznam resursov
terraform state list

# Podrobnosti specifičnega resursa
terraform state show [RESOURCE_NAME]
``` 