# Epileptic

The **Epileptic** application component of the *Communications Platform MVP*
plateau, delivered by the **Setup** work package.

Epileptic is a Go service built on [Gin](https://github.com/gin-gonic/gin). It
realizes the Communication Service and hosts the Content Management function,
which owns the Content data object.

## Model mapping

| Model element                        | Kind                | Code                              |
| ------------------------------------ | ------------------- | --------------------------------- |
| Epileptic                            | ApplicationComponent | this repository (`cmd/server`)    |
| go/gin                               | SystemSoftware      | Gin router in `internal/server`   |
| Communication Service                | ApplicationService  | HTTP API (`internal/server`)      |
| Content Management                   | ApplicationFunction | `internal/content`                |
| Content                              | DataObject          | `content` table (`migrations/`)   |

## Requirements

- Go 1.25+
- Postgres 16+ (declared as the `epileptic-db` service in `.plot/package.json`)

## Configuration

Configuration is read from the environment (see `.env.example`):

| Variable                 | Default                                                                        | Purpose                        |
| ------------------------ | ------------------------------------------------------------------------------ | ------------------------------ |
| `EPILEPTIC_HTTP_ADDR`    | `:8080`                                                                        | Address the service listens on |
| `EPILEPTIC_DATABASE_URL` | `postgres://epileptic:epileptic@localhost:5432/epileptic?sslmode=disable`      | Content data store             |
| `EPILEPTIC_ENV`          | `debug`                                                                        | Gin run mode                   |

## Running

Locally with a Postgres instance:

```sh
make tidy
make run
```

Or the full stack (app + Postgres) with Docker:

```sh
make up
```

Migrations are embedded in the binary and applied on startup.

## HTTP API

| Method   | Path           | Description                  |
| -------- | -------------- | ---------------------------- |
| `GET`    | `/healthz`     | Liveness probe               |
| `GET`    | `/readyz`      | Readiness probe (pings DB)   |
| `POST`   | `/content`     | Create a Content record      |
| `GET`    | `/content`     | List Content records         |
| `GET`    | `/content/:id` | Read a Content record        |
| `PUT`    | `/content/:id` | Update a Content record      |
| `DELETE` | `/content/:id` | Delete a Content record      |

Example:

```sh
curl -X POST localhost:8080/content \
  -H 'Content-Type: application/json' \
  -d '{"title":"Welcome","body":"Hello from Epileptic"}'
```

## Development

```sh
make test    # go test ./...
make lint    # go vet ./...
make fmt     # gofmt
```
