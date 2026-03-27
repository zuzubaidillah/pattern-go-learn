# Go API Boilerplate - Modular Monolith

This project is a high-performance, production-ready REST API boilerplate written in Go. It implements a **Modular Monolith** pattern combined with **Clean Architecture**.

## 🚀 Tech Stack

- **Framework**: [Gin](https://gin-gonic.com/) for high-throughput HTTP routing and middleware.
- **Database**: MySQL via native `database/sql` driver with robust transaction mappings and separate **Read** and **Write** connection pools.
- **Caching**: [Redis](https://redis.io/) via `go-redis` (`cache-aside` strategy).
- **Logging**: [Zap](https://github.com/uber-go/zap) for fast structural logging.
  - **Slack Integration**: Custom integrated Zap core that automatically formats and pushes logs over a certain threshold to a Slack webhook. 
- **Containerization**: Docker Compose configured natively for OrbStack/Docker Desktop to run your local SQL and Redis engines.

---

## 🏗️ Architecture Design (Clean Architecture)

The codebase strictly enforces the dependency rule: `Outer Layer -> Inner Layer`.

1. **Transport / HTTP**: Gin Handlers, routing, and incoming/outgoing payload DTO validation (`internal/modules/*/transport/http`).
2. **Service**: Core business logic and use cases. Coordinates cache invalidations and calls to repositories (`internal/modules/*/service`).
3. **Repository**: Concrete SQL queries mapping database structures directly to memory Domain items (`internal/modules/*/repository`).
4. **Domain**: Business Entities, App Contracts, and Interface abstractions (`internal/modules/*/domain`).

---

## 🛠️ Project Setup

### 1. Prerequisites

- Go 1.20+
- [OrbStack](https://orbstack.dev/) or Docker Desktop (for local DB/Redis spinning)

### 2. Configuration & Environment

Duplicate the example configuration file:
```bash
cp .env.example .env
```
Modify `.env` to suit your requirements (e.g., filling out your personal `SLACK_WEBHOOK_URL`).

### 3. Spin up Infrastructure

Run the Docker containers for MySQL and Redis in detached mode:
```bash
docker compose up -d
```
> **Note:** Included `scripts/init.sql` runs strictly on the first boot of the MySQL container to auto-create the `yourapp_write` and `yourapp_read` databases along with the `users` testing tables.

### 4. Running the Application

**Standard Build/Run:**
```bash
go run cmd/api/main.go
```

**Development Mode (Hot-Reloading):**
To achieve reactive/automatic code reloading upon any file changes, this project utilizes [Air](https://github.com/air-verse/air).
If you do not have Air installed globally, install it first:
```bash
go install github.com/air-verse/air@latest
```
Then, start the application with:
```bash
$(go env GOPATH)/bin/air
```
*(Catatan: Jika Anda hanya mengetik `air` dan muncul program bahasa R, berarti ada bentrok nama program di Mac Anda. Selalu gunakan `$(go env GOPATH)/bin/air`)*

The server will boot and display the **DB Listen** pipeline health checks, starting correctly on `localhost:8080`.

---

## ✨ Demo Module: User Endpoints

A fully built-out `User` namespace is included to serve as a reference implementation for Clean Architecture:
- `GET /api/v1/users/:id`: Fetches a user (Prioritizes Redis Cache -> Hits Read DB -> Set Cache -> Return).
- `POST /api/v1/users`: Creates a new user (Validates Email -> Commits to Write DB -> Returns mapped User).
- `PUT /api/v1/users/:id`: Updates a user profile (Updates Write DB -> Invalidates Cache).
