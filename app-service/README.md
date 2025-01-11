# App Service

Mikrostoritev za upravljanje z aplikacijami in njihovimi metapodatki.

## Zahteve

- Go 1.21+
- PostgreSQL
- Docker (opcijsko)
- Kubernetes (opcijsko)

## Namestitev

1. Klonirajte repozitorij in se premaknite v direktorij:
```bash
cd app-service
```

2. Kopirajte `.env.example` v `.env` in nastavite ustrezne vrednosti:
```bash
cp .env.example .env
```

3. Zaženite migracijske skripte za podatkovno bazo:
```bash
go run cmd/migrate/main.go up
```

## Razvoj

### Lokalno okolje

1. Namestite odvisnosti:
```bash
go mod download
```

2. Zaženite aplikacijo:
```bash
go run cmd/server/main.go
```

### Docker

1. Zgradite Docker image:
```bash
docker build -t app-service .
```

2. Zaženite container:
```bash
docker run -d \
  --name app-service \
  --env-file .env \
  -p 8080:8080 \
  app-service
```

## Testiranje

Za zagon testov uporabite:

```bash
# Enota testi
go test ./...

# Integracijski testi z Docker okoljem
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

## Kubernetes

Za namestitev v Kubernetes okolje uporabite manifeste v `k8s/` mapi. Pred namestitvijo ustvarite Secret z ustreznimi vrednostmi za okolje:

```bash
kubectl create secret generic app-service-secrets \
  --from-literal=DB_USER=your_db_user \
  --from-literal=DB_PASSWORD=your_db_password \
  --from-literal=DB_NAME=your_db_name \
  --from-literal=DB_HOST=your_db_host
```

## Struktura projekta

- `cmd/` - Vstopne točke aplikacije
- `internal/` - Interna implementacija
  - `handler/` - HTTP/gRPC handlers
  - `repository/` - Podatkovni dostop
  - `service/` - Poslovna logika
- `migrations/` - SQL migracijske datoteke
- `k8s/` - Kubernetes manifesti
- `docker/` - Docker konfiguracija
- `scripts/` - Pomožne skripte 