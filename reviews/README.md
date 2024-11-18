# Mikrostoritev za ocene in komentarje

*Opomba*: Ta je spisana v rustu, ne v go-ju.

## Poganjanje
1. Poženi bazo podatkov (`docker/compose/docker-compose.yaml`):
```bash
docker compose -f docker/compose/docker-compose.yaml up -d
```
2. Poženi mikrostoritev (`reviews`):
```bash
cargo run
```
