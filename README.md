<div align="center">

# 📊 Logs Management System

### Simulated flight-logs pipeline: producer → Kafka → consumer → MongoDB, served in real time

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Kafka](https://img.shields.io/badge/Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white)](https://kafka.apache.org)
[![MongoDB](https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white)](https://mongodb.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](LICENSE)

</div>

---

## 📖 Overview

A simulated flight-logs pipeline. A **producer** generates fake flight events and seeds the DB; a **consumer** consumes a Kafka topic, buffers writes, and persists to MongoDB while updating flight state; a **server/web app** serves flights, per-flight events via SSE, and aircrafts.

## 🧩 Components

| Component | Role |
|---|---|
| **Producer** | Seeds initial fake data; appends new flights on an interval |
| **Consumer** | Consumes Kafka topic; batched writes to MongoDB; updates flight state (landed/arrived/cancelled) |
| **Server** | Flights table, per-flight event stream via SSE, aircrafts table |
| **Web** | React dashboard |

## 🏗 Architecture

```
logs-mgmt/
├── producer/   # fake log generator
├── consumer/   # Kafka consumer + buffered flush + intelligence
├── server/     # HTTP server (controllers, SSE)
├── db/         # Postgres (sqlc) + MongoDB
├── web/        # React + Vite dashboard
└── docker-compose.yml
```

## 🛠 Tech stack

`Go` · `Kafka` · `MongoDB` · `Postgres (sqlc)` · `Docker` · `React` · `Vite`

## 📦 Quick start

```bash
git clone https://github.com/m3rashid/logs-mgmt.git
cd logs-mgmt
docker compose up
```

## 📄 License

MIT © [m3rashid](https://github.com/m3rashid)
