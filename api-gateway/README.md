# API Gateway

API Gateway služi kot vstopna točka za našo mikroservisno arhitekturo. Upravlja z usmerjanjem zahtev, avtentikacijo in avtorizacijo.

## Zahteve

- Go 1.21 ali novejši
- Make
- Docker (opcijsko)

## Namestitev

1. Klonirajte repozitorij

2. Kopirajte `.env.example` v `.env` in nastavite ustrezne vrednosti:
```bash
cp .env.example .env
```

3. Namestite odvisnosti:
```bash
go mod download
```

## Zagon

### Lokalno razvijanje
```bash
go run src/main.go
```

### Docker
1. Zgradite Docker image:
```bash
docker build -t api-gateway .
```

2. Zaženite container:
```bash
docker run -d \
  --name api-gateway \
  --env-file .env \
  -p 3000:3000 \
  api-gateway
```

## Razvoj

### Generiranje protokol buffers
```bash
make proto
```

## Konfiguracija

Storitev uporablja naslednje okoljske spremenljivke:

- `MINIAUTH_ADDRESS`: Naslov authentication servisa (privzeto: "localhost:50051")
- `PORT`: Vrata na katerih teče strežnik (privzeto: 3000)

## Kubernetes

Za namestitev v Kubernetes okolje uporabite Kubernetes manifeste v `terraform/kubernetes/` mapi. Pred namestitvijo ustvarite Secret z ustreznimi vrednostmi za okolje:

```bash
kubectl create secret generic api-gateway-secrets \
  --from-literal=MINIAUTH_ADDRESS=your_auth_address
``` 