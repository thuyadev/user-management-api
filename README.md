# User Management API

A RESTful API built with **Go**, **Gin**, **GORM**, and **MongoDB** following **Clean Architecture** principles.

## Features

- User authentication (login & register) with JWT
- Admin user CRUD with paginated data table response
- Category management (CRUD)
- Product management (CRUD)
- Async user activity logging (MongoDB / NoSQL)
- Laravel-style database seeders
- AI-powered product description & category name suggestions (OpenAI-compatible)
- Unit test coverage

## Project Structure

```
project/
├── controllers/     # HTTP handlers (presentation layer)
├── services/        # Business logic
├── repositories/    # Data access layer
├── middleware/      # Auth & admin guards
├── models/          # Domain models & DTOs
├── routes/          # Route definitions
├── database/        # DB connections, migrations, seeders
├── utils/           # Config, JWT, password, responses
├── postman/         # Postman collection
└── main.go
```

## Prerequisites

- Go 1.22+
- Docker & Docker Compose (for PostgreSQL + MongoDB)

## Quick Start

1. **Start databases**

```bash
docker compose up -d
```

2. **Configure environment**

```bash
cp .env.example .env
```

3. **Install dependencies & run**

```bash
go mod tidy
go run main.go
```

Server runs at `http://localhost:8080`

## Default Admin Credentials

| Field    | Value              |
|----------|--------------------|
| Email    | admin@example.com  |
| Password | password123        |

## API Endpoints

### Public
| Method | Endpoint                | Description        |
|--------|-------------------------|--------------------|
| GET    | /health                 | Health check       |
| POST   | /api/v1/auth/login      | Login              |
| POST   | /api/v1/auth/register   | Register           |

### Authenticated
| Method | Endpoint                | Description        |
|--------|-------------------------|--------------------|
| GET    | /api/v1/auth/me         | Current user info  |

### Admin (requires admin role)
| Method | Endpoint                              | Description              |
|--------|---------------------------------------|--------------------------|
| GET    | /api/v1/admin/users                   | List users (paginated)   |
| POST   | /api/v1/admin/users                   | Create user              |
| GET    | /api/v1/admin/users/:id               | Get user                 |
| PUT    | /api/v1/admin/users/:id               | Update user              |
| DELETE | /api/v1/admin/users/:id               | Delete user              |
| GET    | /api/v1/admin/categories              | List categories          |
| POST   | /api/v1/admin/categories              | Create category          |
| GET    | /api/v1/admin/categories/:id          | Get category             |
| PUT    | /api/v1/admin/categories/:id          | Update category          |
| DELETE | /api/v1/admin/categories/:id          | Delete category          |
| POST   | /api/v1/admin/categories/ai/suggest   | AI category suggestion   |
| GET    | /api/v1/admin/products                | List products            |
| POST   | /api/v1/admin/products                | Create product           |
| GET    | /api/v1/admin/products/:id            | Get product              |
| PUT    | /api/v1/admin/products/:id            | Update product           |
| DELETE | /api/v1/admin/products/:id            | Delete product           |
| POST   | /api/v1/admin/products/ai/description | AI product description   |
| GET    | /api/v1/admin/logs                    | List activity logs       |

## Running Tests

```bash
go test ./... -v
```

## Postman Collection

Import `postman/User_Management_API.postman_collection.json` into Postman.

Run **Login (Admin)** first — it auto-saves the JWT token for other requests.

## AI Integration

Set in `.env`:

```
AI_ENABLED=true
AI_API_KEY=your-openai-api-key
```

When disabled, AI endpoints return sensible fallback responses.

## API Security

### API key (required by default)

Every request under `/api/v1` must include your API key header:

```
X-API-Key: your-api-key-from-env
```

Configure in `.env`:

| Variable | Description |
|----------|-------------|
| `API_KEY_REQUIRED` | `true` to enforce (default) |
| `API_KEY` | Secret key — use a long random string in production |
| `API_KEY_HEADER` | Header name (default: `X-API-Key`) |

`/health` stays public for load balancers and does not require the key.

Postman: the collection sends `X-API-Key` automatically via the `api_key` variable.

### CORS (browser protection)

Browsers sending an `Origin` header are only allowed if the origin is listed in `CORS_ALLOWED_ORIGINS` (comma-separated). Empty = block all cross-origin browser requests.

Postman, curl, and server-to-server calls are not restricted by CORS.

### Security headers

All responses include `X-Content-Type-Options`, `X-Frame-Options`, `Cache-Control: no-store`, and related headers.

### Recommended practices for production

| Practice | Why |
|----------|-----|
| **HTTPS only** | Encrypts API key, JWT, and passwords in transit |
| **Strong `API_KEY` & `JWT_SECRET`** | Use 32+ random bytes; rotate periodically |
| **Rate limiting** | Mitigate brute-force on `/auth/login` (add at reverse proxy or with a limiter middleware) |
| **Short JWT expiry + refresh** | Limits damage if a token is stolen |
| **Never commit `.env`** | Keep secrets in environment / secret manager |
| **Validate & sanitize input** | Already handled via Gin binding; keep DTOs strict |
| **Least privilege** | Admin routes require `admin` role (already enforced) |
| **Audit logs** | User CRUD events go to MongoDB asynchronously |
| **Database TLS** | Enable `DB_SSLMODE=require` and MongoDB TLS in production |
| **Disable `API_KEY_REQUIRED` only in local dev** | Never turn off in production |

Disable API key locally only if needed:

```
API_KEY_REQUIRED=false
```

## Architecture Notes

- **RDBMS (PostgreSQL)**: Users, Categories, Products via GORM
- **NoSQL (MongoDB)**: User activity logs written asynchronously via goroutines
- **Clean Architecture**: Controllers → Services → Repositories with interface-based DI
