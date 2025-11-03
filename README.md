# Titanbay Take-Home — Backend Service

This repository contains a compact backend service for managing **funds**, **investors**, and **investments**. The service exposes a small REST API backed by PostgreSQL and is designed to run locally with minimal setup.

---

## Contents

* [Technology Stack](#technology-stack)
* [Getting Started (Docker)](#getting-started-docker)
* [Getting Started (Local)](#getting-started-local)
* [Configuration](#configuration)
* [API Overview](#api-overview)
* [Data Models](#data-models)
* [Examples](#examples)
* [Testing](#testing)
* [Assumptions & Design Decisions](#assumptions--design-decisions)
* [Project Structure](#project-structure)
* [Troubleshooting](#troubleshooting)
* [Notes & Future Work](#notes--future-work)

---

## Technology Stack

* **Language:** Go
* **HTTP:** Echo framework
* **Database:** PostgreSQL (via GORM)
* **Identifiers:** UUIDs for primary keys
* **Monetary values:** `shopspring/decimal` for precise currency handling
* **Containerization:** Docker & Docker Compose for local development

---

## Getting Started (Docker)

### Prerequisites

* Docker Desktop or Docker Engine with Docker Compose v2+

### Run

```bash
cd dev
docker compose up --build
```

**Services**

* API: `http://localhost:1323`
* PostgreSQL: `database:5432` (internal to the compose network)

**Environment (from `docker-compose.yml`)**

* `PORT=:1323`
* `DB_URL=host=database user=tb_user password=tb_pass dbname=tb_tbdb port=5432 sslmode=disable`

Database schema is created/updated automatically via GORM migrations at startup.

---

## Getting Started (Local)

### Prerequisites

* Go ≥ 1.21
* PostgreSQL ≥ 14

### Steps

1. **Create database and user**

   ```sql
   CREATE USER tb_user WITH PASSWORD 'tb_pass';
   CREATE DATABASE tb_tbdb OWNER tb_user;
   ```

2. **Set environment variables**

   ```bash
   export PORT=":1323"
   export DB_URL="host=localhost user=tb_user password=tb_pass dbname=tb_tbdb port=5432 sslmode=disable"
   ```

3. **Run the service**

   ```bash
   go mod download
   go run main.go
   ```

   The server will listen on `http://localhost:1323`.

---

## Configuration

| Variable      | Description                                 | Example                                   |
| ------------- | ------------------------------------------- | ----------------------------------------- |
| `PORT`        | Listen address                              | `:1323`                                   |
| `DB_URL`      | PostgreSQL DSN (GORM format)                | `host=localhost user=... sslmode=disable` |

> In Docker, these are provided by `dev/docker-compose.yml`.

---

## API Overview

The service implements eight REST endpoints as described in the assignment brief and accompanying API specification.

| Method | Path                          | Description                        |
| -----: | ----------------------------- | ---------------------------------- |
|    GET | `/funds`                      | List funds                         |
|   POST | `/funds`                      | Create a fund                      |
|    PUT | `/funds`                      | Update a fund                      |
|    GET | `/funds/:fund_id`             | Retrieve a fund by UUID            |
|    GET | `/investors`                  | List investors                     |
|   POST | `/investors`                  | Create an investor                 |
|    GET | `/funds/:fund_id/investments` | List investments                   |
|   POST | `/funds/:fund_id/investments` | Create an investment               |

**Error Handling**
Validation and database errors are returned as JSON with appropriate HTTP status codes. In production, internal details would be redacted.

---

## Data Models

**Fund**

* `id` (uuid)
* `name` (string)
* `vintage_year` (int)
* `target_size_usd` (decimal)
* `status` (enum: `Fundraising | Investing | Closed`)
* `created_at` (string: date-time)

**Investor**

* `id` (uuid)
* `name` (string)
* `investor_type` (string; e.g., `Individual`, `Institutional`)
* `email` (string)
* `created_at` (string: date-time)

**Investment**

* `id` (uuid)
* `fund_id` (uuid, FK -> Fund)
* `investor_id` (uuid, FK -> Investor)
* `amount_usd` (decimal)
* `investment_date` (date)
* `created_at` (string: date-time)

---

## Examples

### Create a Fund

```bash
curl -X POST http://localhost:1323/funds \
  -H 'Content-Type: application/json' \
  -d '{
        "name": "TB Secondary II",
        "vintage_year": 2021,
        "target_size_usd": "250000000",
        "status": "Fundraising"
      }'
```

### List Funds

```bash
curl http://localhost:1323/funds
```

### Create an Investor

```bash
curl -X POST http://localhost:1323/investors \
  -H 'Content-Type: application/json' \
  -d '{
        "name": "Alex Rivera",
        "investor_type": "Individual",
        "email": "alex@example.com"
      }'
```

### Create an Investment (Commitment)

```bash
curl -X POST http://localhost:1323/investments \
  -H 'Content-Type: application/json' \
  -d '{
        "fund_id": "UUID-OF-FUND",
        "investor_id": "UUID-OF-INVESTOR",
        "commitment_usd": "5000000",
        "commitment_date": "2024-03-15"
      }'
```

---

## Testing

Run unit tests for handlers:

```bash
go test ./handlers
```

Tests cover routing, handler behavior, and validation at the HTTP boundary.

---

## Assumptions & Design Decisions

* **Scope fidelity:** The implementation adheres closely to the brief and spec, avoiding additional features (e.g., authentication, advanced filtering) to maintain clarity and focus.
* **Echo + GORM:** Selected for mature ecosystems, succinct APIs, and rapid iteration.
* **UUIDs:** Used as stable, non-sequential identifiers suitable for client exposure.
* **Decimal for money:** Ensures precise handling of monetary values and avoids floating-point errors.
* **Validation:** Model-level validation methods invoked by handlers; invalid requests return HTTP 400 with human-readable messages.
* **Error responses:** Consistent JSON envelope. `{ "error": "..." }`. Internal errors should be redacted in production.
* **Developer experience:** Container-first workflow enables one-command local startup without installing Go/PostgreSQL.
* **Migrations:** Auto-migrate is enabled for reviewer convenience; production would use versioned migrations.

---

## Project Structure

```
.
├─ main.go
├─ handlers/               # HTTP handlers and tests
├─ models/                 # Domain models and validation
├─ database/               # DB connection, DB interface, mock DB, and GORM setup
└─ dev/                    # Dockerfile and docker-compose for local run
```

---

## Troubleshooting

* **“PORT not set”**
  Ensure `PORT=":1323"` is exported (or use Docker which sets this automatically).

* **Database connection failures**
  Confirm `DB_URL` (host, credentials, database name). With Docker, verify the DB container is healthy:

  ```bash
  docker compose ps
  ```

* **Validation errors**
  Review the JSON error response to identify missing or invalid fields.

---

## Notes & Future Work

* The implementation was time-boxed; priority was given to correctness, validation, and smooth local setup.
* Potential enhancements:

  * Structured logging and request IDs
  * Versioned migrations
  * Pagination, filtering, and sorting on list endpoints
  * Serving an OpenAPI document directly from the service
  * CI for tests and linting

---

**Thank you for reviewing this submission.**
