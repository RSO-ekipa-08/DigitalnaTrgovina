# Terraform Infrastruktura

Terraform konfiguracija za postavitev infrastrukture v oblaku.

## Zahteve

- Terraform 1.0+
- Azure CLI
- kubectl

## Struktura projekta

- `kubernetes/` - Konfiguracija za Kubernetes cluster
- `main.tf` - Glavna Terraform konfiguracija
- `variables.tf` - Definicije spremenljivk
- `outputs.tf` - Izhodne vrednosti
- `providers.tf` - Konfiguracija ponudnikov
- `terraform.tfvars.example` - Primer konfiguracijskih vrednosti

## Namestitev

1. Namestite Azure CLI in se prijavite:
```bash
az login
```

2. Kopirajte `terraform.tfvars.example` v `terraform.tfvars` in nastavite ustrezne vrednosti:
```bash
cp terraform.tfvars.example terraform.tfvars
```

3. Inicializirajte Terraform:
```bash
terraform init
```

## Uporaba

### Pregled sprememb

Pred apliciranjem sprememb preverite plan:

```bash
terraform plan
```

### Apliciranje sprememb

Za postavitev infrastrukture:

```bash
terraform apply
```

### Brisanje infrastrukture

Za odstranitev vse infrastrukture:

```bash
terraform destroy
```

## Konfiguracija

Glavne spremenljivke v `terraform.tfvars`:

- `resource_group_name` - Ime Azure resource group
- `location` - Azure regija
- `kubernetes_version` - Verzija Kubernetes clustra
- `node_count` - Število worker vozlišč
- `vm_size` - Velikost Azure VM-jev

## Izhodni podatki

Po uspešni postavitvi so na voljo naslednji izhodni podatki:

- Kubernetes cluster credentials
- Resource group podatki
- Endpoint naslovi

Za prikaz izhodnih podatkov:

```bash
terraform output
```

## Varnost

- Občutljivi podatki (gesla, ključi) naj bodo shranjeni v `terraform.tfvars`
- Ne commitajte `terraform.tfvars` ali `*.tfstate` datotek v git
- Uporabljajte state backend s primerno zaščito (npr. Azure Storage Account)

## Vzdrževanje

### Posodobitev modulov

Za posodobitev Terraform modulov:

```bash
terraform init -upgrade
```

### Upravljanje state datoteke

Za prikaz trenutnega stanja:

```bash
terraform show
```

Za seznam resursov:

```bash
terraform state list
``` 