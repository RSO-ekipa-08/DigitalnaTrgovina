# Dokumentacija mikrostoritev

## Pregled arhitekture

Sistem je sestavljen iz naslednjih mikrostoritev:

- **App Service** (Go) - upravljanje z aplikacijami
- **Reviews Service** (Rust) - upravljanje z ocenami
- **Authentication Service** (Go) - avtentikacija preko Auth0
- **Payment Service** (Bun/TypeScript) - integracija s plačilnim sistemom
- **API Gateway** (Go) - vstopna točka za vse zahteve

## Komunikacija med storitvami

### Protokoli

- **gRPC** - primarna komunikacija med storitvami
- **GraphQL** - podpora za Reviews Service
- **REST** - API Gateway ponuja REST vmesnik za zunanje odjemalce
- **RabbitMQ** - asinhrona komunikacija za plačila

### Varnost

- JWT žetoni za avtentikacijo
- CORS zaščita
- Rate limiting

## Mikrostoritve

### App Service
Glavna storitev za upravljanje z aplikacijami v trgovini.

**Funkcionalnosti:**
- CRUD operacije za aplikacije
- Iskanje in filtriranje aplikacij
- Upravljanje kategorij
- Procesiranje prenosov
- Integracija s plačilnim sistemom

**Tehnologije:**
- Go
- PostgreSQL
- gRPC/Connect
- Minio za shranjevanje datotek

**API:**
- REST API preko API Gateway

### Reviews Service
Upravlja z ocenami in komentarji aplikacij.

**Funkcionalnosti:**
- Dodajanje ocen in komentarjev
- Moderiranje vsebine
- Statistika ocen
- Multi-tenant podpora

**Tehnologije:**
- Rust
- PostgreSQL
- GraphQL
- gRPC

**API:**
- GraphQL API
- gRPC API definiran v `reviews.proto`

### Authentication Service
Skrbi za avtentikacijo uporabnikov preko Auth0.

**Funkcionalnosti:**
- Prijava/odjava uporabnikov
- Preverjanje JWT žetonov
- Upravljanje s sejami

**Tehnologije:**
- Go
- Auth0
- gRPC

**API:**
- gRPC API definiran v `auth.proto`
- REST API preko API Gateway

### Payment Service
Procesira plačila preko Stripe.

**Funkcionalnosti:**
- Kreiranje plačil
- Webhook za Stripe dogodke
- Zgodovina plačil

**Tehnologije:**
- Bun/TypeScript
- RabbitMQ
- Stripe API

**API:**
- RabbitMQ sporočila
- REST webhook

### API Gateway
Vstopna točka za vse zunanje zahteve.

**Funkcionalnosti:**
- Usmerjanje zahtev
- Avtentikacija
- Rate limiting
- CORS
- Dokumentacija API-ja

**Tehnologije:**
- Go
- gRPC-Gateway

**API:**
- REST API za zunanje odjemalce
- https://rso-ekipa-08.github.io/DigitalnaTrgovina/

## Namestitev in zagon

### Zahteve
- Docker & Docker Compose
- Go 1.22+
- Rust 1.75+
- Node.js 20+
- PostgreSQL 16
- RabbitMQ

### Koraki

1. Kloniranje repozitorija:
```bash
git clone <repo-url>
```

2. Nastavitev okolja:
```bash
cp .env.example .env
```

3. Zagon s kompozicijo:
```bash
docker-compose up -d
```

4. Zagon posameznih storitev:

**App Service:**
```bash
cd app-service
make run
```

**Reviews Service:**
```bash
cd reviews
cargo run
```

**Auth Service:**
```bash
cd authentication
go run .
```

**Payment Service:**
```bash
cd payment_v2
bun run src/index.ts
```

**API Gateway:**
```bash
cd api-gateway
go run .
```

## Razvoj

### API spremembe

1. Posodobitev proto definicij
2. Generiranje novih stub-ov
3. Implementacija sprememb
4. Testi
5. Uveljavitev sprememb

### Dodajanje nove storitve

1. Nova mapa z ustrezno strukturo
2. Proto definicija API-ja
3. Implementacija storitve
4. Dodajanje v Docker kompozicijo
5. Integracija z API Gateway
6. Testiranje


## Produkcijska namestitev

1. Build Docker slik:
```bash
docker-compose build
```

2. Push slik v register:
```bash
docker-compose push
```

3. Namestitev na Kubernetes:
```bash
kubectl apply -f k8s/
```

## Spremljanje

Na Azure je na voljo:
- Prometheus metrike
- Grafana dashboard

## Dodatne informacije

### Glej tudi

- Semantic versioning
- Conventional commits
- Proto style guide
- Go style conventions
- Rust style guide
