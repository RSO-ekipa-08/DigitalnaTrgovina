# Payment Service

Mikrostoritev za obdelavo plačil preko Stripe API-ja z uporabo RabbitMQ za komunikacijo.

## Zahteve

- Node.js
- Bun
- RabbitMQ
- Stripe račun in API ključ

## Namestitev

1. Namestite odvisnosti:
```bash
bun install
```

2. Kopirajte `.env.example` v `.env` in nastavite ustrezne vrednosti:
```bash
cp .env.example .env
```

3. Nastavite Stripe API ključ in druge spremenljivke v `.env` datoteki.

## Zagon

### Lokalno
```bash
bun run index.ts
```

### Docker
1. Zgradite Docker image:
```bash
docker build -t payment-service .
```

2. Zaženite container:
```bash
docker run -d \
  --name payment-service \
  --env-file .env \
  --network host \
  payment-service
```

Opomba: Za povezavo z RabbitMQ v Dockerju lahko uporabite omrežje `host` (kot zgoraj) ali pa definirate Docker omrežje in uporabite ime RabbitMQ container-ja kot hostname.

## Uporaba

Storitev sprejema sporočila preko RabbitMQ vrste `payment_requests` v formatu:

```json
{
    "amount": 1000,
    "currency": "eur",
    "productName": "Ime produkta",
    "quantity": 1
}
```

Odgovori so poslani v vrsto `payment_responses` v formatu:

```json
{
    "checkoutUrl": "https://checkout.stripe.com/..."
}
```

## Kubernetes

Za namestitev v Kubernetes okolje uporabite Kubernetes manifeste v `k8s/` mapi. Pred namestitvijo ustvarite Secret z ustreznimi vrednostmi za okolje:

```bash
kubectl create secret generic payment-service-secrets \
  --from-literal=STRIPE_SECRET_KEY=your_stripe_key \
  --from-literal=SUCCESS_URL=your_success_url \
  --from-literal=CANCEL_URL=your_cancel_url \
  --from-literal=RABBITMQ_URL=your_rabbitmq_url
```
