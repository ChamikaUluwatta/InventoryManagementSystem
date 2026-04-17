# Inventory Management System

An open-source inventory management system for small-scale hardware and grocery shops.

The goal of this project is to provide a simple, practical, and extensible solution for day-to-day stock management.

## Current Feature Status

- [x] Product management (View, Add, Edit)
- [] Location management (View, Add, Edit)
- [] Inventory management (View, Add, Edit)
- [] Category management (View, Add, Edit)
- [] Company management (View, Add, Edit)
- [] Role Based management (for View,Add,Edit)
- [] Sales satistic (view)
- [] Integrated E-Commerce site
## Tech Stack

- Backend: Go (net/http), PostgreSQL
- Frontend: React + TypeScript + Vite
- Containerization: Docker + Docker Compose

## Project Structure

- `Backend/` - Go API server, domain modules, migrations
- `Frontend/` - React web application
- `docker-compose.yml` - Multi-service local deployment

## Environment Variables

### Backend

Create `Backend/.env` when running manually:

```env
DB_HOST=postgres://postgres:postgres@localhost:5432/inventory?sslmode=disable
SERVER_PORT=8080
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

### Frontend

Create `Frontend/.env` when running manually:

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
```

## Running Migrations

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

- Backend API: http://localhost:8080
- Frontend: http://localhost:5173
- PostgreSQL: localhost:5432 (or your remapped host port)

### Manual Approach

1. Start PostgreSQL (local or Docker).
2. Apply migrations.
3. Start backend.
4. Start frontend.

Backend (from `Backend/`):

```bash
go mod download
go run ./cmd
```

Frontend (from `Frontend/`):

```bash
npm install
npm run dev
```

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
