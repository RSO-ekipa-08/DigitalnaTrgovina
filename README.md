# Digitalna Trgovina - Mikrostoritve

Digitalna Trgovina je sistem mikrostoritev za upravljanje digitalne trgovine z aplikacijami. Sestavljen je iz več ključnih mikrostoritev:

## Mikrostoritve

### 1. App Service
Osnovna storitev za upravljanje z aplikacijami v trgovini.
- Implementirano v Go
- gRPC/Connect API
- PostgreSQL podatkovna baza
- Supabase za shranjevanje datotek

### 2. Reviews Service
Storitev za upravljanje z ocenami in komentarji aplikacij.
- Implementirano v Rust
- gRPC API + GraphQL
- PostgreSQL podatkovna baza

### 3. Authentication Service
Storitev za avtentikacijo uporabnikov preko Auth0.
- Implementirano v Go
- gRPC API
- Integracija z Auth0

### 4. Payment Service
Storitev za procesiranje plačil.
- Implementirano v JavaScript (Bun runtime)
- RabbitMQ za komunikacijo
- Integracija s Stripe

### 5. API Gateway
Centralna vstopna točka za vse mikrostoritve.
- Implementirano v Go
- Usmerja zahteve na ustrezne mikrostoritve
- Upravlja avtentikacijo in avtorizacijo

## Zahteve sistema

### Razvoj
- Go 1.22+
- Rust 1.75+
- Node.js/Bun
- Docker in Docker Compose
- Make
- Protocol Buffers (protoc)

### Podatkovne baze
- PostgreSQL 16
- RabbitMQ
- Supabase/MinIO (za shranjevanje datotek)

### Infrastruktura
- Auth0 račun
- Stripe račun (za plačila)
- Kubernetes (opcijsko za produkcijo)

## Namestitev

1. Kloniraj repozitorij:
```bash
git clone https://github.com/RSO-ekipa-08/DigitalnaTrgovina.git
cd DigitalnaTrgovina
```

2. Ustvari potrebne .env datoteke za vsako storitev po vzorcu .env.example

3. Zaženi podatkovne baze:
```bash
docker compose up -d postgres rabbitmq minio
```

4. Zaženi posamezne storitve:

**App Service:**
```bash
cd app-service
make migrate
make run
```

**Reviews Service:**
```bash
cd reviews
cargo run
```

**Authentication Service:**
```bash
cd authentication
go run src/main.go
```

**Payment Service:**
```bash
cd payment_v2
bun run src/index.ts
```

**API Gateway:**
```bash
cd api-gateway
go run src/main.go
```

## Razvoj

### Generiranje protokol buffer datotek
```bash
make proto # v vsaki storitvi
```

## API Dokumentacija

https://rso-ekipa-08.github.io/DigitalnaTrgovina/

## Arhitektura

Sistem uporablja mikroservisno arhitekturo s sledečimi značilnostmi:
- gRPC za sinhrono komunikacijo med storitvami
- RabbitMQ za asinhrono komunikacijo (plačila)
- API Gateway za usmerjanje zahtev
- Auth0 za avtentikacijo
- PostgreSQL za shranjevanje podatkov
- Supabase/MinIO za shranjevanje datotek

## Konfiguracija

Vse storitve uporabljajo okoljske spremenljivke za konfiguracijo. Glej .env.example datoteke v vsaki storitvi za potrebne nastavitve.
