# DevBoard

A small project/task/time-tracking API written in Go (Gin + GORM + SQLite + JWT auth).

## Stack

- **Gin** — HTTP router/framework
- **GORM** + **SQLite** — ORM and embedded database (no external DB server needed)
- **golang-jwt** — JWT-based auth
- **bcrypt** — password hashing

## Setup

1. Install Go 1.24+ (or whatever your local toolchain supports — `go.mod` targets a recent version; lower it if your installed Go is older).
2. From the project root:
   ```bash
   go mod tidy
   go run ./cmd/main.go
   ```
3. The server starts on the port in `.env` (default `8080`). A `devboard.db` SQLite file is created automatically on first run and tables are migrated on startup.

`.env` (already present):
```
DEVBOARD_PORT=8080
DEVBOARD_DATABASE=devboard.db
DEVBOARD_JWT_SECRET=change-this-development-secret
```
Change `DEVBOARD_JWT_SECRET` before deploying anywhere real.

## Quick smoke test

```bash
# Register
curl -X POST localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada","email":"ada@example.com","password":"password123"}'

# Login (grab the "token" from the response)
curl -X POST localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ada@example.com","password":"password123"}'

# Use the token for everything else
TOKEN="paste-token-here"

curl -X POST localhost:8080/api/projects \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"DevBoard","description":"Build the thing"}'

curl localhost:8080/api/reports/dashboard -H "Authorization: Bearer $TOKEN"
```

## API

All routes below except `/api/auth/*` and `/health` require `Authorization: Bearer <token>`.

**Auth**
- `POST /api/auth/register` — `{name, email, password}`
- `POST /api/auth/login` — `{email, password}` → `{token, user}`

**Profile**
- `GET /api/profile`
- `PUT /api/profile` — `{name, email}`
- `DELETE /api/profile`

**Projects**
- `POST /api/projects` — `{name, description}`
- `GET /api/projects`
- `GET /api/projects/:project_id`
- `PUT /api/projects/:project_id` — `{name, description, status}`
- `DELETE /api/projects/:project_id`

**Tasks**
- `POST /api/projects/:project_id/tasks` — `{title, description, priority}`
- `GET /api/projects/:project_id/tasks`
- `GET /api/tasks/:task_id`
- `PUT /api/tasks/:task_id` — `{title, description, status, priority}`
- `DELETE /api/tasks/:task_id`

**Time tracking**
- `POST /api/time/start` — `{project_id, task_id}`
- `POST /api/time/stop/:id`
- `GET /api/time`
- `GET /api/projects/:project_id/time`
- `GET /api/tasks/:task_id/time`

**Reports**
- `GET /api/reports/projects/:project_id/time` — total time logged on a project
- `GET /api/reports/tasks/productivity` — time spent per task, ordered by most time
- `GET /api/reports/activity` — overall counts + total tracked time for the user
- `GET /api/reports/dashboard` — activity summary + task-status breakdown + 5 most recent time entries

## What to extend first

1. **Change-password route** — `UserService.UpdatePassword` already exists but has no handler/route wired up. Good first PR: add a `PUT /api/profile/password` endpoint.
2. **Pagination** on `GET /api/projects`, `/api/tasks`, `/api/time` — right now they return everything for the user, fine for a prototype but will get slow with real data.
3. **Refresh tokens** — JWTs currently expire in 24h with no refresh flow; add `POST /api/auth/refresh`.
4. **Team support** — everything is single-user-owns-everything right now (`user_id` on every row). A real "board" tool usually wants shared projects with roles (the `Role` field on `User` already hints at this but isn't enforced anywhere yet).
5. **Tests** — there are currently none; the service layer (`internal/services`) is the easiest place to start since it's plain functions over a `*gorm.DB`.
