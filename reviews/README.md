# Mikrostoritev za ocene in komentarj
# Mikrostoritev za ocene in komentarje

Mikrostoritev, implementirana v programskem jeziku Rust, ki omogoča upravljanje z ocenami in komentarji aplikacij v trgovini.
Uporablja gRPC protokol za komunikacijo in PostgreSQL za shranjevanje podatkov.

## Funkcionalnosti

- Dodajanje ocen in komentarjev
- Branje posameznih in seznama ocen
- Posodabljanje obstoječih ocen
- Brisanje ocen
- Moderiranje komentarjev
- Podpora za več najemnikov (multi-tenancy)
- Izračun povprečnih ocen in statistik

## Tehnične zahteve

- Rust 1.75 ali novejši
- PostgreSQL 16
- Docker in Docker Compose (opcijsko)
- protoc (Protocol Buffers prevajalnik)

## Namestitev in zagon

### 1. Priprava okolja

```bash
# Namestitev potrebnih orodij
cargo install sqlx-cli
```

### 2. Nastavitev podatkovne baze

```bash
# Zagon PostgreSQL z Docker Compose
docker compose -f docker/docker-compose.yaml up -d

# Ali uporaba obstoječe PostgreSQL instance:
# Nastavi DATABASE_URL v .env datoteki:
DATABASE_URL=postgres://reviews_user:reviews_password@localhost:5432/reviews_db
```

### 3. Migracije podatkovne baze

```bash
# Izvedba migracij
sqlx migrate run --source db/migrations
```

### 4. Zagon storitve

```bash
# Razvojno okolje
SQLX_OFFLINE=true cargo run

# Produkcija
SQLX_OFFLINE=true cargo run --release
```

## Uporaba Docker okolja

```bash
# Build Docker Image
docker build -t reviews-service -f docker/Dockerfile .

# Zagon vsebnika
docker run -p 50051:50051 \
  -e DATABASE_URL=postgres://reviews_user:reviews_password@db:5432/reviews_db \
  reviews-service
```

## API dokumentacija

gRPC vmesnik je definiran v `proto/reviews.proto`
HHTML dokumentacija je na voljo v `docs/index.html`.

## Razvoj

```bash
# Generiranje novih migracij
sqlx migrate add <ime_migracije>

# Posodobitev sqlx-data.json za offline preverjanje
cargo sqlx prepare

# Preverjanje kode
cargo check
cargo test
cargo clippy
```
