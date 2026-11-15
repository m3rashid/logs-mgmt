# Logs Management System

Simulated flight logs pipeline.

## Components

- **Producer** — generates fake flight logs, seeds DB, appends flights on interval
- **Consumer** — consumes Kafka topic, buffers writes to MongoDB, updates flight state
- **Server** — flights table, per-flight events via SSE, aircrafts table

## Tech

Go · Kafka · MongoDB · Docker

## License

MIT
