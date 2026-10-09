# Epileptic

The **Epileptic** application component of the **Communications Platform MVP**
plateau. It realizes the **Communication Service** and exposes the **Content
Management** function over the **Content** data object.

This repository is the deliverable of the **Setup** work package: it stands up
the toolchain, the service skeleton, and its backing Postgres service.

## Toolchain

| Concern    | Choice                                   |
| ---------- | ---------------------------------------- |
| Language   | Go (see `go.mod`)                         |
| HTTP       | [Gin](https://github.com/gin-gonic/gin)   |
| Persistence| Postgres via [pgx](https://github.com/jackc/pgx) |

The toolchain and the declared services live in
[`.plot/package.json`](.plot/package.json).

## Model mapping

| Model element            | Artefact                                              |
| ------------------------ | ----------------------------------------------------- |
| Communication Service    | HTTP router + endpoints in `internal/server`          |
| Epileptic component      | `cmd/server` and the `internal/...` packages          |
| Content Management       | `internal/content` (handler + repository)             |
| Content (data object)    | `content` table in `internal/database`                |
| go/gin (system software) | `go.mod` dependencies                                 |
| Postgres service         | `docker-compose.yml` + `internal/database`            |

## Configuration

Copy `.env.example` to `.env` and adjust values. Configuration is read from the
environment:

- `PORT` – HTTP port (default `8080`).
- `DATABASE_URL` – Postgres connection string. When unset, it is assembled from
  `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`, and `PGSSLMODE`.

## Running

```sh
docker compose up -d postgres   # start the local Postgres service
make run                        # start the Communication Service
```

Build and test:

```sh
make build
make test
```

## API

| Method | Path            | Description                        |
| ------ | --------------- | ---------------------------------- |
| GET    | `/healthz`      | Liveness probe                     |
| GET    | `/readyz`       | Readiness probe (pings Postgres)   |
| GET    | `/content`      | List Content                       |
| POST   | `/content`      | Create Content                     |
| GET    | `/content/:id`  | Read a single Content record       |
| PUT    | `/content/:id`  | Update a Content record            |
| DELETE | `/content/:id`  | Delete a Content record            |
