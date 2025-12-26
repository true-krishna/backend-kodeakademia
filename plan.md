**Project Overview**
- **Purpose**: Build a Udemy-like learning platform with owner-managed course CRUD, learner features (browse, buy, mark complete), Google OAuth, PostgreSQL, REST API in Go using Clean Architecture.
- **Primary actors**: **Owner** (platform owner creates courses), **Learner** (regular user browses and buys courses).

**Phases & High-Level TODOs**
- **Project Skeleton & Dependencies**: Initialize `go.mod`, repository layout (`cmd/`, `internal/`, `pkg/`, `migrations/`), add `docker-compose.yml` stub, `Makefile`, linter (golangci-lint), and migration tool (golang-migrate). Acceptance: repo builds and linters run.
- **Database Schema & Migrations**: Create SQL migrations for `users`, `courses`, `purchases`, `course_progress`, `roles`. Acceptance: migrations run and tables exist.
- **Domain Models & Repository Interfaces**: Define domain entities and repository interfaces in `internal/domain/`.
- **Postgres Repository Implementations**: Implement repository interfaces in `internal/adapters/postgres/` with unit tests/mocks.
- **Google OAuth Authentication**: Implement Google OAuth login, user provisioning, and JWT/session issuance; add auth middleware.
- **Usecases / Business Logic**: Implement usecases in `internal/usecase/` (list/get/buy/complete course, owner CRUD, stats) with unit tests.
- **HTTP Handlers and Routing**: Create REST handlers and routes in `internal/transport/http/` with DTOs and validation.
- **Validation & Error Handling**: Add request/response DTOs, input validation (e.g., `go-playground/validator`), and consistent error format.
- **Tests (Unit & Integration)**: Unit tests for repos and usecases; integration tests for HTTP handlers against a test Postgres instance.
- **Statistics & Owner Admin**: Owner-only endpoints for course/user statistics and revenue aggregations.
- **CI/CD & Quality Gates**: Add GitHub Actions to run tests, linters, and build; optionally build Docker image.
- **Docs & Deployment**: Provide `README.md`, `.env.example`, and `docker-compose.yml` that runs Postgres + the app.

**Implementation Granularity & Conventions**
- **One responsibility per function**: Functions should do a single thing (e.g., `GetCourseByID`, `CreatePurchase`). Target ~20–100 LOC per function.
- **Layer separation**:
  - **Handler**: HTTP parsing, auth, call usecase, format response.
  - **Validation/DTO**: Request structs + validation tags and mapping to domain.
  - **Usecase/Service**: Business logic; uses repository interfaces only.
  - **Repository**: DB access, SQL queries, transactions.
  - **Mapper**: Domain ↔ transport conversions.
- **File layout recommendations**:
  - `internal/domain/entity/*.go`
  - `internal/domain/repository/*.go`
  - `internal/adapters/postgres/*.go`
  - `internal/usecase/*.go`
  - `internal/transport/http/{handlers,middleware,dto}.go`
  - `migrations/0001_init.sql`, etc.

**Test Strategy & Small-step Workflow**
- **Repository tests**: Use a test DB or DB mocks; cover success and common DB failures.
- **Usecase tests**: Mock repositories; cover happy path, validation, and repository failures.
- **Handler tests**: Use `httptest` or integration tests against a test DB and a test JWT provider.
- **Per-function micro-steps (recommended)**:
  1. Add interface/signature and compile.
  2. Add unit-test file with failing tests.
  3. Implement function body.
  4. Run `go test` and fix until green.
  5. Add integration test if cross-layer behavior matters.

**Concrete Example: "Buy a Course" feature (micro-steps)**
- **A — DTO & Route**: Add `internal/transport/http/dto/purchase.go` with `BuyCourseRequest` + validation; add route `POST /courses/{id}/buy`.
- **B — Repo Interface**: Add `CreatePurchase(ctx, purchase) error` to `internal/domain/repository/purchase_repo.go` and a stub in `internal/adapters/postgres/purchase_repo.go`.
- **C — Usecase**: Implement `BuyCourse(ctx, userID, courseID)` in `internal/usecase/purchase.go` (check course exists, ensure not already purchased, record purchase in a transaction). Write unit tests mocking repos.
- **D — Handler**: Implement HTTP handler to parse auth+path, call `BuyCourse`, and return 201 or appropriate error code. Add integration test hitting the endpoint with test DB and JWT.
- **E — Acceptance**: Integration test verifies a row in `purchases` and that the user can access the course content.

**Example local commands**
```bash
# start Postgres (using docker-compose)
docker-compose up -d postgres

# run migrations (golang-migrate example)
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# run unit tests
go test ./...
```

**Next Steps**
- I can scaffold the project layout now (create `go.mod`, `cmd/server/main.go`, `docker-compose.yml`, `migrations/` stub and basic `README.md`). This implements the first TODO (Project skeleton & dependencies). Reply with "Scaffold" to proceed or pick another TODO ID to start.

**Files**
- This plan is saved as `plan.md` in the repository root of the workspace.

**Technology Stack**
- **Language**: Go (1.20+ recommended)
- **Web framework**: `Echo` (`github.com/labstack/echo/v4`) — request routing, middleware, and HTTP handlers.
- **Database**: PostgreSQL (containerized via `docker-compose`).
- **DB driver / query tooling**: `pgx` (`github.com/jackc/pgx/v5`) for low-level access; optionally use `sqlc` for type-safe SQL code generation or `sqlx` for convenience.
- **Migrations**: `golang-migrate/migrate` (`github.com/golang-migrate/migrate`) to apply SQL migrations from `migrations/`.
- **OAuth & Auth**:
  - OAuth2 (OAuth 2.0) client: `golang.org/x/oauth2` — implement standard OAuth 2.0 flows for Google login.
  - JWT handling: `github.com/golang-jwt/jwt/v5` (or v4 depending on compatibility)
- **Validation**: `github.com/go-playground/validator/v10` for request validation.
- **Logging**: `go.uber.org/zap` for structured logging (or `sirupsen/logrus` if preferred).
- **Configuration**: `github.com/kelseyhightower/envconfig` or `github.com/spf13/viper` for env/config management.
- **Testing**:
  - Unit testing: Go `testing` package + `github.com/stretchr/testify` for assertions and mocks.
  - HTTP tests: `net/http/httptest` and Echo's testing helpers.
- **Linting / Formatting / Quality**:
  - `golangci-lint` for static analysis
  - `gofmt` / `gofumpt` for formatting
  - `go vet`
- **CI/CD**: GitHub Actions to run tests, linters, and build Docker images.
- **Containerization / Local dev**: `docker`, `docker-compose` (Postgres + Adminer dev UI)
- **Dev utilities**:
  - DB client: `psql` or Adminer
  - Migrations CLI: `migrate`
  - Secrets: `.env` for local dev, recommend using a secrets manager for production

**Notes / rationale**
- `Echo` is chosen per request for routing and middleware simplicity and good performance.
- `pgx` is recommended for Postgres performance and feature completeness; `sqlc` pairs nicely if you prefer writing SQL and generating types.
- `golang-migrate` keeps migrations plain SQL and easy to run in CI and containers.

**Local Database Options (manual + Docker)**

If you prefer to create the database manually while still using Docker for the Postgres image, use one of the two options below.

Option A — Docker (single container) + manual DB creation (no docker-compose)
- Start Postgres container:
```bash
docker run -d --name kodeakademia-postgres -p 5432:5432 \
  -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=postgres postgres:15
```
- Create user and database manually with `psql` (inside container):
```bash
docker exec -it kodeakademia-postgres psql -U postgres -c "CREATE USER app_user WITH PASSWORD 'secret';"
docker exec -it kodeakademia-postgres psql -U postgres -c "CREATE DATABASE app_db OWNER app_user;"
```
- Apply migrations from your host using `migrate` or `psql` (example using `migrate`):
```bash
export DB_USER=app_user DB_PASSWORD=secret DB_NAME=app_db DB_HOST=localhost DB_PORT=5432
migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up
```
- Seed the DB:
```bash
psql "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME" -f migrations/seeds/seed_courses.sql
```

Option B — Manual Postgres binary on host (no Docker)
- Install Postgres locally and create DB/user directly with `createdb` / `createuser` or via `psql`.
- Apply migrations and seed as shown above, adjusting `DB_HOST` to `localhost` and credentials accordingly.

Notes:
- If you use Option A, make sure the port `5432` is free on your host or map a different host port (e.g., `-p 5433:5432`) and update `DB_PORT` in env.
- The existing `docker-compose.yml` remains available as an alternative if you later prefer a compose workflow.



**Feature-by-feature Implementation Plan (one-by-one)**

Note: features are developed sequentially. Each feature below is separate and self-contained — do not combine work across features until the previous one passes its acceptance criteria.
-Feature 0 — Project scaffold & Health Check (TDD)
- Goal: Scaffold the Go project, add a minimal Echo server, and a `GET /health` endpoint to verify runtime wiring.
- Endpoint(s): `GET /health` (returns 200 and basic JSON {"status":"ok"}).
- Files & symbols to add:
  - `go.mod` with module name and dependencies (Echo, zap, pgx placeholders).
  - `cmd/server/main.go` — minimal Echo server, env loading, and route registration.
  - `internal/transport/http/health.go` — health handler and tests.
  - `.gitignore`, a minimal `Makefile` target `run` and `test`.
  - `README.md` short note on how to run the server and test the health endpoint.
- Testing & TDD specifics:
  - Follow Red-Green-Refactor: write a failing unit test for the `health` handler (expecting HTTP 200), then implement the handler to make the test pass, then refactor.
  - Add a simple integration smoke test that starts the server (httptest) and calls `GET /health`.
- Acceptance criteria: `go test ./...` passes for the health handler tests; running the server (`make run` or `go run ./cmd/server`) and requesting `http://localhost:8080/health` returns HTTP 200 and JSON with `status: ok`.

-Feature 1 — User Login (Google OAuth 2.0)
- Goal: Allow users to authenticate via Google using OAuth 2.0, provision user records, and receive a JWT for protected endpoints.
- Endpoint(s): `GET /auth/google/login` (redirect), `GET /auth/google/callback` (callback)
- Files & symbols to add:
  - `internal/adapters/auth/google.go` — handles OAuth token exchange and fetches Google profile.
  - `internal/usecase/auth.go` — `LoginOrProvisionUser(ctx, googleProfile) (domain.User, tokenString, error)`.
  - `internal/domain/repository/user_repo.go` — ensure `CreateOrUpdateFromOAuth(ctx, profile) (User, error)` exists.
  - `internal/transport/http/routes_auth.go` — route registration.
  - `internal/transport/http/middleware/auth.go` — JWT validation and context injection.
- DB changes: `users` table (email unique, google_sub nullable/unique, name, avatar_url, role)
- Tests:
  - Unit: `LoginOrProvisionUser` with mocked user repo (happy path, existing user, repo error).
  - Integration: simulate OAuth callback flow using a mock HTTP server for Google token/profile endpoints and a test Postgres instance.
- Acceptance criteria: After callback, client receives a JWT; protected endpoints accept the JWT and return user-specific responses.

Feature 2 — Course List (public)
- Goal: Provide a paginated, searchable listing of courses for all users.
- Endpoint(s): `GET /courses` (query params: `page`, `size`, `q`, `category`)
- Files & symbols to add:
  - `internal/domain/repository/course_repo.go` — `List(ctx, params) ([]Course, Pagination, error)`.
  - `internal/adapters/postgres/course_repo.go` — implement SQL with LIMIT/OFFSET and optional WHERE clauses.
  - `internal/usecase/course_list.go` — `ListCourses(ctx, params)` returning DTOs and pagination metadata.
  - `internal/transport/http/dto/course_list.go` — request/query parsing and response DTO.
  - `internal/transport/http/handlers/courses.go` — handler for `GET /courses`.
- DB considerations: index on `title`, `category`, and `created_at` for performance.
- Tests:
  - Unit: usecase tests mocking course repo (pagination boundaries, filters, empty results).
  - Integration: seed test DB with sample courses and call `GET /courses` asserting pagination and content.
- Acceptance criteria: `GET /courses` returns status 200, correct `total`, `page`, `size`, and `items` matching filters.

Feature 3 — Course Detail (public)
- Goal: Return full course details including owner info and aggregated stats (enrollments, completions, avg rating if later added).
- Endpoint(s): `GET /courses/{id}`
- Files & symbols to add:
  - `internal/domain/repository/course_repo.go` — `GetByID(ctx, id) (Course, error)` and `GetStats(ctx, id) (CourseStats, error)`.
  - `internal/adapters/postgres/course_repo.go` — SQL to fetch course and compute aggregates (COUNT of purchases, COUNT of completions).
  - `internal/usecase/course_detail.go` — `GetCourseDetail(ctx, id)` combining course, owner, and stats.
  - `internal/transport/http/dto/course_detail.go` — response DTO.
  - `internal/transport/http/handlers/courses.go` — handler for `GET /courses/{id}`.
- DB considerations: use efficient aggregates (materialized views later if needed).
- Tests:
  - Unit: usecase tests mocking repos (not found, success, repo error).
  - Integration: seed DB and verify `GET /courses/{id}` returns expected numbers for enrollments and completions.
- Acceptance criteria: `GET /courses/{id}` returns 200 with `id`, `title`, `owner{name,id}`, `price`, `enrollments`, `completions`.

Development rules for sequential work
- Complete each feature end-to-end (unit tests green, integration test passing) before starting the next.
- Keep each feature's commits focused and titled like `feat(auth): google oauth login`.
- For each feature, follow the micro-step workflow: define interfaces → add failing tests → implement → run tests → add integration test.

**Development Process (TDD)**
- We will follow Test-Driven Development (TDD) using the Red-Green-Refactor pattern for all new code and bug fixes.
- Red (write failing test): write a focused unit test that expresses the desired behavior or interface before implementing the code. Tests must be small and fast.
- Green (make it pass): implement the minimal amount of production code to make the test pass. Prefer simple, correct implementations over premature optimization.
- Refactor: clean up code, remove duplication, improve names, and maintain test coverage. Run tests after each change to ensure no regressions.
- Practical TDD rules for the team:
  - Write tests first for each usecase/function. For HTTP handlers you may write handler-level tests (httptest) first, then mock underlying usecases/repositories for unit tests.
  - Keep unit tests isolated: mock network, DB, and external OAuth providers. Use integration tests only to validate cross-layer behavior with a real/test DB.
  - Each commit that adds functionality must include the tests that demonstrate it (commit message prefix `test:` when tests are added/updated).
  - Use small, incremental commits following the Red→Green→Refactor cycle; avoid large commits that mix many changes.
  - Aim for clear, deterministic tests (no reliance on wall-clock time, randomness, or external systems without mocks).
  - Maintain a fast test suite for rapid Red-Green cycles — keep unit tests under ~200ms where practical and funnel slower scenarios to integration suites.
  - CI must run the full test suite (unit + integration) on PRs; local development should encourage fast unit runs and selective integration runs.
  - Use table-driven tests for multiple cases and descriptive test names for clarity.

**Testing & Coverage Targets**
- Aim for high coverage on critical business logic (usecases/repositories). A realistic target is 80%+ for core packages, but prioritize meaningful tests over raw coverage percent.


