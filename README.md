<div align="center">

# 📊 Logs Management System

**Simulated flight-logs pipeline: producer → Kafka → consumer → MongoDB, served in real time.**

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

---

## 🧩 Components

| Component | Role |
|---|---|
| **Producer** | Seeds fake data; appends flights on an interval |
| **Consumer** | Consumes Kafka topic; batched writes to MongoDB; updates flight state (landed/arrived/cancelled) |
| **Server/Web** | Flights table, per-flight event stream via SSE, aircrafts table |

## 🛠 Tech

`Go` · `Kafka` · `MongoDB` · `Postgres (sqlc)` · `Docker`

## 📦 Quick start

```bash
git clone https://github.com/m3rashid/logs-mgmt.git
cd logs-mgmt
docker compose up
```

## 📄 License

MIT
