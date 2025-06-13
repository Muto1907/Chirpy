# Chirpy

Chirpy is a minimal Twitter-like micro-blogging API written in Go. It showcases a clean server architecture with JWT authentication, a PostgreSWL backend and SQLC-generated data access layer.

> ✨  Spin it up, hit the API, and start *chirping!*

---
## Motivation

This project was undertaken to deepen my understanding of http server design. The server authenticates Users using a combination of JWT and refresh Tokens.

---
## Features

|   | Feature                                                             |
| - | ------------------------------------------------------------------- |
| ✅ | **Pure standard‑library HTTP server**—no external frameworks        |
| ✅ | **JWT access tokens** (HS256, issuer *chirpy‑access*, 1 h)          |
| ✅ | Hex‑encoded **refresh tokens** (256 random bytes → 512‑char string) |
| ✅ | Passwords stored with **bcrypt** (current cost = *MinCost*)         |
| ✅ | PostgreSQL with **sqlc**‑generated, type‑safe queries               |
| ✅ | *Bad‑word* filter ("kerfuffle", "sharbert", "fornax" → `****`)      |
| ✅ | Static file hosting under **/app** (ideal for a SPA front‑end)      |
| ✅ | Basic HTML metrics page at **/admin/metrics** (counts static hits)  |
| ✅ | Dev‑only **/admin/reset** endpoint to wipe DB + metrics             |
| ✅ | **Polka** webhook → auto‑upgrade users to *Chirpy Red*              |

---

## Quick start

```bash
# 1 – Clone & enter
git clone https://github.com/Muto1907/Chirpy && cd Chirpy

# 2 – Copy & edit env file
cp .env.example .env

# 3 – Start Postgres (local, Docker, whatever you like)

# 4 – Run the API
go run .
```

Health‑check:

```bash
curl http://localhost:8080/api/healthz   # → "OK"
```

---

## Environment variables

| Name         | Required | Example                                                          | Purpose                                    |
| ------------ | -------- | ---------------------------------------------------------------- | ------------------------------------------ |
| `DB_URL`     | ✔        | `postgres://chirpy:chirpy@localhost:5432/chirpy?sslmode=disable` | PostgreSQL DSN                             |
| `SECRET_KEY` | ✔        | `super-long-hmac-secret`                                         | HMAC secret for JWTs                       |
| `POLKA_KEY`  | ✔        | `pka_live_123`                                                   | Shared secret for Polka webhook            |
| `PLATFORM`   | ✖        | `dev`                                                            | Enables **/admin/reset** when set to `dev` |
| `PORT`       | ✖        | `8080`                                                           | HTTP listen port (default 8080)            |

> The server automatically loads an `.env` file if it exists.

---

## Security design

| Layer                | What we do                                                                                                                                                             | Why it matters                                                              |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------- |
| **Password storage** | Hash with *bcrypt* (Go’s `golang.org/x/crypto/bcrypt`, cost = `bcrypt.MinCost`).                                                                                       | Straightforward, industry‑standard defence against credential leaks.        |
| **Access token**     | JWT v5, algorithm **HS256**, issuer `chirpy-access`, subject = user UUID. Validity configurable (default 1 h). Signed with `SECRET_KEY`.                               | Self‑contained, verifiable credentials—no DB lookup needed for every call.  |
| **Refresh token**    | 256 cryptographically random bytes → hex‑encoded string (512 chars). Stored in DB with expiry (30 days) and **revocation flag**.                                       | Long‑lived token that can be invalidated server‑side without touching JWTs. |
| **Headers**          | All protected routes expect `Authorization: Bearer <token>` where `<token>` is **either** a JWT (access routes) **or** a refresh‑token string (refresh/revoke routes). | Keeps the wire format consistent and simple.                                |
| **Token rotation**   | Hitting `/api/refresh` issues a new JWT but keeps the same refresh token until it expires or you explicitly revoke at `/api/revoke`.                                   | Minimises DB chatter yet allows immediate logout across devices.            |

---

## API reference

All endpoints are **JSON over HTTP**. Always send

```
Content-Type: application/json
```

Responses come back with the same header. Authentication, when required, is via

```
Authorization: Bearer <token>
```
### Endpoint specification

Below is a **complete** list of routes exposed by Chirpy.

| Method     | Path                  | Purpose                            |
| ---------- | --------------------- | ---------------------------------- |
| **POST**   | `/api/users`          | Register new user                  |
| **PUT**    | `/api/users`          | Update email/password              |
| **POST**   | `/api/login`          | Log in & obtain tokens             |
| **POST**   | `/api/refresh`        | Exchange refresh token for new JWT |
| **POST**   | `/api/revoke`         | Revoke a refresh token             |
| **POST**   | `/api/chirps`         | Create chirp                       |
| **GET**    | `/api/chirps`         | List chirps (filter/sort)          |
| **GET**    | `/api/chirps/{id}`    | Fetch one chirp                    |
| **DELETE** | `/api/chirps/{id}`    | Delete own chirp                   |
| **POST**   | `/api/polka/webhooks` | Polka upgrade webhook              |
| **GET**    | `/api/healthz`        | Health‑check (plain text)          |
| **GET**    | `/admin/metrics`      | HTML metrics page                  |
| **POST**   | `/admin/reset`        | Dev‑only: reset DB & metrics       |
| **GET**    | `/app/*`              | Serve static web assets            |

---

### User and Password resource

User:

```json
{
    "id": "qweqweklsdfdkfl-wsqeasd",
    "created_at": 2012-04-23T18:25:43.511Z,
    "updated_at": 2012-04-23T18:25:43.511Z,
    "email": "Mudi@example.com",
    "is_chirpy_red": false
}
```

Password:

```json
{
    "password":"wewewe2",
    "email":"Mudi@example.com"
}
```

### POST /api/users

Creates a new user

Request Body:
```json
{
    "password":"wewewe2",
    "email":"Mudi@example.com"
}
```

Response Body:
```json
{
    "id": "qweqweklsdfdkfl-wsqeasd",
    "created_at": 2012-04-23T18:25:43.511Z,
    "updated_at": 2012-04-23T18:25:43.511Z,
    "email": "Mudi@example.com",
    "is_chirpy_red": false
}
```

### PUT /api/users

Updates Email and Password of a user. Access Token must be included in Request Header in the following form:
"Authorization": "Bearer TOKENHERE"

Request Body:
```json
{
    "password":"wewewe2",
    "email":"Mudi@example.com"
}
```

Response Body:
```json
{
    "id": "qweqweklsdfdkfl-wsqeasd",
    "created_at": 2012-04-23T18:25:43.511Z,
    "updated_at": 2012-04-23T18:25:43.511Z,
    "email": "Mudi@example.com",
    "is_chirpy_red": false
}
```


