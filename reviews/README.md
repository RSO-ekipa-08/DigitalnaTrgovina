# Mikrostoritev za ocene in komentarje

*Opomba*: Ta je spisana v rustu, ne v go-ju.

## Poganjanje
1. Poženi bazo podatkov (`docker/compose/docker-compose.yaml`):
```bash
docker compose -f docker/docker-compose.yaml up -d
```
2. Migracije baze podatkov (`cargo install sqlx-cli`):
```bash
sqlx migrate run --source db/migrations
```
2. Poženi mikrostoritev (`reviews`):
```bash
SQLX_OFFLINE=true cargo run
```
SQLX_OFFLINE omogoča gradnjo z offline preverjanjem s podatki iz mape `.sqlx`.
Podatke v `json` obliki dobimo z ukazom
```bash
cargo sqlx prepare
```
