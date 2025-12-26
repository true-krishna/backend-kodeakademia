# Kodeakademia — Backend

This repository contains the backend for a simple learning platform (a Udemy-like app). The project uses Go, follows Clean Architecture principles, and is developed using TDD (Red-Green-Refactor).

This README explains the repository layout, how to run the app locally, how to run migrations and seeds, and the development process (including TDD rules and how features are broken down).

## Status

- Feature 0 (project scaffold + health check): implemented.
- Feature 1 (Google OAuth 2.0 login): planned in `plan.md` and ready to implement.

See `plan.md` for the full feature-by-feature plan and development TODOs.

## Technology Stack

- Language: Go (1.20+)
- Web framework: Echo (`github.com/labstack/echo/v4`)
- Database: PostgreSQL
- DB driver: `pgx` recommended (`github.com/jackc/pgx/v5`)
- Migrations: `golang-migrate/migrate`
- OAuth2: `golang.org/x/oauth2` (Google OAuth 2.0)
- JWT: `github.com/golang-jwt/jwt/v5`
- Validation: `github.com/go-playground/validator/v10`
- Logging: `go.uber.org/zap`
- Testing: Go `testing`, `github.com/stretchr/testify`

## Repository layout (important paths)

- `cmd/server/` — application entrypoint (`main.go`)
- `internal/` — private application code
	- `internal/domain/` — domain entities and repository interfaces
	- `internal/usecase/` — business logic (usecases)
	- `internal/adapters/` — adapters (Postgres, OAuth, etc.)
	- `internal/transport/httptransport/` — HTTP handlers and middleware
- `migrations/` — SQL migrations and seed files
- `plan.md` — development plan, feature list, and TDD rules

## Environment variables

Use `.env` for local development (the repo includes `.env.example`). Important vars:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `DATABASE_URL` / `MIGRATE_DATABASE_URL` — convenience URL for migration tools
- `APP_BASE_URL` — e.g. `http://localhost:8080`
- `OAUTH_REDIRECT_URI` — e.g. `http://localhost:8080/auth/google/callback`
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`
- `JWT_SECRET` — generate a secure random secret (`openssl rand -hex 32`)
- `TEST_DB_*` — optional test DB creds for integration tests

Note: If your DB password contains special characters (for example `@`), you must URL-encode the password when you set `DATABASE_URL`/`MIGRATE_DATABASE_URL`.

## Getting started (development)

Prereqs:

- Go 1.20+
- Docker (optional for Postgres)
- `psql` or Adminer (Adminer is included in `docker-compose.yml` if you want a UI)

1. Copy `.env.example` to `.env` and fill values (or ensure your environment provides these variables).

2. Start Postgres (choose one):

- Option A — Docker (quick start using `docker run`):

```bash
docker run -d --name kodeakademia-postgres -p 5432:5432 \
	-e POSTGRES_PASSWORD=secret -e POSTGRES_USER=postgres postgres:15

docker exec -it kodeakademia-postgres psql -U postgres -c "CREATE USER app_user WITH PASSWORD 'secret';"
docker exec -it kodeakademia-postgres psql -U postgres -c "CREATE DATABASE app_db OWNER app_user;"
```

- Option B — Docker Compose (predefined):

```bash
docker-compose up -d postgres adminer
```

3. Run migrations (example with `migrate`):

```bash
export DB_USER=app_user DB_PASSWORD=secret DB_NAME=app_db DB_HOST=localhost DB_PORT=5432
migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up
```

4. (Optional) Seed development data:

```bash
psql "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME" -f migrations/seeds/seed_courses.sql
```

5. Run the server (Feature 0 scaffold):

```bash
go mod download
make run
# or
go run ./cmd/server
```

6. Verify health endpoint:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Running tests

Unit tests (fast):

```bash
make test
```

Integration tests:

- Integration tests require a running Postgres instance. Use the `TEST_DB_*` vars in `.env` to point to a test database and run migrations/seeds against it before executing integration tests.

## Development Process (TDD)

All development follows TDD (Red-Green-Refactor):

1. Red — write a failing unit test that specifies the required behavior.
2. Green — implement the smallest amount of code to make the test pass.
3. Refactor — improve design and clean up while keeping tests green.

Practical rules:

- Unit tests must be fast and isolated (mock DB and external services).
- Integration tests validate end-to-end behaviors with a real/test DB.
- Each commit adding functionality must include tests demonstrating it.

## OAuth 2.0 (Google) notes

To enable Google OAuth 2.0 login:

1. Create an OAuth 2.0 Client ID in Google Cloud Console.
2. Add the redirect URI (for local development): `http://localhost:8080/auth/google/callback` and set `OAUTH_REDIRECT_URI` accordingly.
3. Set `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` in `.env`.

During development you can either use real Google credentials (easy but requires browser redirect) or mock the token/profile endpoints in tests using `httptest.Server`.

## Migrations & Seeding

- Migrations live in `migrations/` as plain SQL files (`0001_init.up.sql` / `.down.sql`).
- Seeds are in `migrations/seeds/`.
- Use `golang-migrate` CLI or Docker image to apply migrations in CI or locally.

## Linting & Formatting

- Use `gofmt` / `gofumpt` to format code.
- Use `golangci-lint` to run static analysis checks.

## Feature plan and workflow

The full feature plan (sequential features and per-feature micro-steps) is in `plan.md`. Work proceeds feature-by-feature in TDD cycles; Feature 0 is scaffolded and complete.

Feature examples:

- Feature 0 — Project scaffold & health check (done)
- Feature 1 — Google OAuth 2.0 login (planned)
- Feature 2 — Course listing (public endpoint)
- Feature 3 — Course detail (public endpoint)

## Contributing

- Fork and open a pull request with focused changes and tests.
- Follow the Red-Green-Refactor workflow and keep commits small and atomic.
- Run `go test ./...` and `golangci-lint run` before submitting PRs.

## Contact / Questions

If you need me to scaffold additional features or wire OAuth tests, tell me which TODO to start and I will implement the Red step (add failing tests) then the Green (implementation).
