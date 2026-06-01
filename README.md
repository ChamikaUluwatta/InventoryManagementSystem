# Inventory Management System

A multi-tenant inventory management platform with real-time stock tracking, role-based access control, and per-tenant data isolation via PostgreSQL Row-Level Security.

![Architecture](docs/architecture-diagram.png)

## Tech Stack

| Layer        | Technology                                                            |
| ------------ | --------------------------------------------------------------------- |
| Frontend     | React 19 · TypeScript · Vite · TailwindCSS 4 · shadcn/ui              |
| Backend API  | Go 1.25 · chi · pgx · golang-jwt · go-redis                           |
| Auth Service | Java 21 · Spring Boot 4 · Spring Security · Nimbus JOSE+JWT · Flyway   |
| Database     | PostgreSQL 15 (Row-Level Security per `tenant_id`)                     |
| Cache        | Redis 7 (refresh tokens · JWT version keys)                           |
| Proxy        | Nginx (rate limiting · SSL · SPA static serving)                      |
| Infra        | Docker Compose · Cloudflare · DigitalOcean Droplet                    |

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+, Java 21+, Node 20+ (for native dev)

### Development

```bash
# Full stack via Docker (Postgres + Redis in containers, services on host)
make docker/up

# Or run each service natively
make run/backend
make run/frontend
make run/authservice
```

### Database Migrations

```bash
make migrate/up          # apply all
make migrate/down/1      # rollback one step
make migrate/create name=add_x   # create new migration pair
```

## API Documentation

OpenAPI specs:
- [Backend API](docs/openapi/backend.yaml) (`:8080`)
- [Auth Service API](docs/openapi/authservice.yaml) (`:8081`)

## Deploy

Production deployment is tag-based:

```bash
./deploy.sh v1.0.0
```

Checks out the tag, rebuilds images, and runs `docker compose -f docker-compose.prod.yml up -d`.

## Project Layout

```
AuthService/   # Java Spring Boot auth (JWT signing, JWKS endpoint)
Backend/       # Go inventory API (chi router, RLS middleware, JWKS validation)
Frontend/      # React SPA (Vite, Tailwind, shadcn/ui)
nginx/         # Reverse proxy configs (dev + prod)
db-init/       # PostgreSQL init scripts to create databases(auth_db + inventory)
docs/          # Documentations
```
