# Inventory Management System

## Project Status

![IMPORTANT](https://img.shields.io/badge/IMPORTANT-IN%20DEVELOPMENT-red?style=for-the-badge)


> This project is under active development and is not production-ready.
> Features are still being implemented, and APIs or UI behavior may change without notice.

An open-source inventory management system for small-scale hardware and grocery shops.

The goal of this project is to provide a simple, practical, and extensible solution for day-to-day stock management.

## Current Feature Status

[✓] Implemented
[x] Not implemented yet

- [✓] Product management (View, Add, Edit)
- [✓] Location management (View, Add, Edit)
- [✓] Inventory management (View, Add, Edit)
- [✓] Category management (View, Add, Edit)
- [✓] Company management (View, Add, Edit)
- [✓] Authentication & Authorization (JWT + Redis)
- [x] Role Based management (View, Add, Edit)
- [x] Sales statistic (View)
- [x] Integrated E-Commerce site

## Tech Stack

- **AuthService**: Java 21, Spring Boot 4.0, PostgreSQL, Redis
- **Backend**: Go (net/http), PostgreSQL
- **Frontend**: React + TypeScript + Vite
- **Containerization**: Docker + Docker Compose


## Project Structure

- `AuthService/` - Java Spring Boot authentication service
- `Backend/` - Go API server, domain modules, migrations
- `Frontend/` - React web application

## Authentication & Authorization

The system uses JWT (JSON Web Tokens) with RSA-256 signing for authentication.

### Key Features

- **JWT Access Tokens**: Short-lived (15 min), contain user permissions
- **Refresh Tokens**: Long-lived (30 days), stored in Redis, httpOnly cookies
- **Permission Versioning**: Instant permission revocation across services
- **JWKS Endpoint**: Public key distribution for token verification

### Auth Flow

1. User logs in via `POST /api/v1/auth/login`
2. AuthService validates credentials, issues JWT + refresh token
3. JWT contains `permissions_version` claim for cache invalidation
4. Backend verifies JWT signature using JWKS public key
5. Backend checks Redis for permission version changes
6. If permissions changed, token is rejected (401) → frontend auto-refreshes

### Permission Version Check

When roles or permissions change:
1. AuthService increments `permissions_version` in database
2. AuthService writes new version to Redis (`auth:user:version:{userId}`)
3. Backend rejects tokens with stale version (401 Unauthorized)
4. Frontend automatically refreshes token with updated permissions

## Environment Variables

### AuthService

Create `AuthService/.env`:

```env
DB_URL=jdbc:postgresql://localhost:5432/auth_db
DB_USERNAME=postgres
DB_PASSWORD=postgres
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Backend

Create `Backend/.env`:

```env
DB_HOST=postgres://postgres:postgres@localhost:5432/inventory?sslmode=disable
SERVER_PORT=8080
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
REDIS_HOST=localhost
REDIS_PORT=6379
AUTH_JWKS_URL=http://localhost:8081/api/v1/auth/.well-known/jwks.json
```

### Frontend

Create `Frontend/.env`:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## Database Setup

You can run PostgreSQL either with Docker or as a local service.

### Option A: Docker PostgreSQL

Use the included compose service:

```bash
docker compose up -d postgres
```

If port 5432 is already in use on your machine, change the host mapping in `docker-compose.yml` from:

```yaml
- "5432:5432"
```

to:

```yaml
- "5433:5432"
```

If you use 5433 on the host for manual backend runs, set:

```env
DB_HOST=postgres://postgres:postgres@localhost:5433/inventory?sslmode=disable
```

### Option B: Local PostgreSQL

Create database and user credentials that match your `DB_HOST` value. Example:

```sql
CREATE DATABASE inventory;
CREATE DATABASE auth_db;
```

## Running Migrations

### AuthService Migrations

AuthService uses Flyway for automatic migrations. Migrations run on startup.

Migration files are in:
- `AuthService/src/main/resources/db/migration/`

### Backend Migrations

Migration files are in:
- `Backend/internal/database/migrations/`

Example with golang-migrate CLI (from `Backend/`):

```bash
migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/inventory?sslmode=disable" up
```

To roll back one step:

```bash
migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/inventory?sslmode=disable" down 1
```

## Deployment and Run

### Docker Approach (Recommended)

From repository root:

```bash
docker compose up --build
```

Services:

- AuthService: http://localhost:8081
- Backend API: http://localhost:8080
- Frontend: http://localhost:5173
- PostgreSQL: localhost:5432
- Redis: localhost:6379

### Manual Approach

1. Start PostgreSQL (local or Docker).
2. Start Redis.
3. Apply backend migrations.
4. Start AuthService.
5. Start Backend.
6. Start Frontend.

**Redis** (required for AuthService):

```bash
redis-server
```

or with Docker:

```bash
docker run -d -p 6379:6379 redis:7-alpine
```

**AuthService** (from `AuthService/`):

```bash
./mvnw spring-boot:run
```

**Backend** (from `Backend/`):

```bash
go mod download
go run ./cmd
```

**Frontend** (from `Frontend/`):

```bash
npm install
npm run dev
```

## API Endpoints

Available in docs/openapi

## Optional Seed Mode

The backend supports a seed flag:

```bash
go run ./cmd -seed
```

## Contributing

Contributions, suggestions, and issue reports are welcome.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request
