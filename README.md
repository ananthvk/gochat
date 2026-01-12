# GoChat

A real-time chat application with Go backend and React frontend.

## Quick Start

### Docker (recommended)

```bash
# Start PostgreSQL
docker volume create pg_data
docker run -d --name gochat-db \
  -e POSTGRES_PASSWORD=dev \
  -p 5432:5432 \
  -v pg_data:/var/lib/postgresql/data \
  postgres:16-alpine

# Set database URL
export GOCHAT_DB_DSN="postgres://postgres:dev@localhost:5432/postgres?sslmode=disable"

# Run migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -database $GOCHAT_DB_DSN -path internal/database/migrations up

# Run the app
GOCHAT_ENV=development go run ./cmd/gochat
```

The app will be available at http://localhost:8000


You can also run the docker image from docker hub, 

```
docker run --network host --env-file .env ananthvk0/gochat:0.0.4
```

### Development

Frontend:

```bash
cd frontend
npm install
npm run dev
```

Backend (from project root):

```bash
export GOCHAT_DB_DSN="postgres://postgres:dev@localhost:5432/postgres?sslmode=disable"
export GOCHAT_ENV=development
go run ./cmd/gochat
```

### Docker build

```bash
docker build --build-arg VITE_API_BASE_URL=http://localhost:8000/api/v1 -t gochat .
docker run --env-file .env -p 8000:8000 gochat
```

## Migrations

```bash
# Create new migration
migrate create -ext sql -dir internal/database/migrations -seq <name>

# Run migrations
migrate -database $GOCHAT_DB_DSN -path internal/database/migrations up

# Rollback
migrate -database $GOCHAT_DB_DSN -path internal/database/migrations down 1
```