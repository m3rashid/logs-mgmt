# Logs Management System

Simulated flight-logs pipeline: a producer generates fake flight events, a Kafka-backed consumer buffers and persists them to MongoDB, and a web app serves flights/events in real time.

## Components

- **Producer** — seeds initial fake data; appends new flights on an interval
- **Consumer** — consumes Kafka topic; batched writes to MongoDB; updates flight state (landed/arrived/cancelled)
- **Server/Web** — flights table, per-flight event stream via SSE, aircrafts table

## Tech

Go · Kafka · MongoDB · Postgres (sqlc) · Docker

## Quick start

```bash
git clone https://github.com/m3rashid/logs-mgmt.git
cd logs-mgmt
docker compose up
```

## License

MIT
