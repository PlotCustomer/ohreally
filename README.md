# Epileptic

The **Epileptic** application component of the **Communications Platform MVP**
(plateau `17425347-0956-40c2-be43-6e6170553581`). It is delivered by the
**Setup** work package.

## What this is

Epileptic is built on **go/gin** and, per the System Design view, it:

- **realises** the `Communication Service` (which in turn realises the
  `epileptic.com` product serving `User`, composed of `Patient` and `Relative`),
- is **composed of** the `Content Management` application function,
- whose `Content Management` function **accesses** the `Content` data object.

This setup therefore provides a runnable Go/Gin server with the Content
Management function exposed over the `Content` data object.

## Running

```sh
go run .
```

The server listens on `0.0.0.0:8080` (the port declared in
`.plot/package.json` → `preview.port`). Override it with `PORT`.

## Endpoints

| Method | Path                | Description                       |
| ------ | ------------------- | --------------------------------- |
| GET    | `/`                 | Service information               |
| GET    | `/health`           | Liveness probe                    |
| GET    | `/api/content`      | List Content                      |
| POST   | `/api/content`      | Create Content                    |
| GET    | `/api/content/:id`  | Read a Content record             |
| PUT    | `/api/content/:id`  | Update a Content record           |
| DELETE | `/api/content/:id`  | Delete a Content record           |

Example:

```sh
curl -X POST http://localhost:8080/api/content \
  -H 'content-type: application/json' \
  -d '{"title":"Welcome","body":"Hello","author":"Setup"}'
```

## Persistence

The `Content` data object is persisted in the default local **postgres**
service declared in `.plot/package.json` when `DATABASE_URL` is set (for
example `postgres://postgres:postgres@localhost:5432/epileptic?sslmode=disable`).
When no database is configured the component falls back to an in-memory store so
the preview can always start.

## Toolchain

Go 1.27.1, declared in `.plot/package.json` and `.mise.toml`. Dependencies are
vendored under `vendor/`, so the server builds without network access.
