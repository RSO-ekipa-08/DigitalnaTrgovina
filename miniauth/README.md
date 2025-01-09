# Auth0 Avtentikacijska Mikrostoritev

Mikrostoritev za upravljanje z avtentikacijo uporabnikov preko platforme Auth0.
Implementira tako gRPC kot HTTP REST API.

## Funkcionalnosti

- Prijava uporabnikov preko Auth0
- Preverjanje veljavnosti žetonov
- Odjava uporabnikov
- Pridobivanje podatkov o uporabniku
- Dvojni API (gRPC + HTTP)

## Tehnične zahteve

- Go 1.22 ali novejši
- Docker
- Make
- Protobuf prevajalnik (protoc)

## Namestitev

1. Prenesite repozitorij:
```bash
git clone <repository-url>
cd authentication
```

2. Ustvarite `.env` datoteko in dodajte Auth0 nastavitve:
```bash
AUTH0_CLIENT_ID=your_client_id
AUTH0_DOMAIN=your_domain.auth0.com
AUTH0_CLIENT_SECRET=your_client_secret
AUTH0_CALLBACK_URL=http://localhost:3000/callback
```

3. Namestite Go dependecies:
```bash
go mod download
```

4. Naredi proto datoteke:
```bash
make proto
```

## Zagon

### Lokalno

```bash
go run src/main.go
```

### Docker

```bash
docker build -t auth0-golang-web-app .
docker run --env-file .env -p 50051:50051 -p 3000:3000 -it auth0-golang-web-app
```

## API Dostopne točke

### HTTP API

- `GET /` - Začetna stran
- `GET /login` - Prijava uporabnika
- `GET /callback` - Auth0 callback
- `GET /user` - Podatki o uporabniku
- `GET /logout` - Odjava uporabnika

### gRPC API

- `Login` - Začetek prijavnega procesa
- `Verify` - Preverjanje avtentikacijskega zahtevka
- `Logout` - Odjava uporabnika
- `VerifyToken` - Preverjanje veljavnosti žetona

## Varnost

Mikrostoritev uporablja:
- OAuth 2.0 / OpenID Connect protokol
- JWT žetone za avtentikacijo
- HTTPS za prenos podatkov (v produkciji)
- Varno shranjevanje občutljivih podatkov v .env datoteki

## Razvoj

Za razvoj novih funkcionalnosti:

1. Posodobite proto definicije v `proto/auth.proto`
2. Generirajte nove proto datoteke z `make proto`
3. Implementirajte novo funkcionalnost
4. Testirajte spremembe
5. Posodobite dokumentacijo
