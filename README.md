# DVBD.com Webshop

Implementation of the **Webshop Setup** WorkPackage of the **Webshop MVP**
plateau. The system is modelled with Archimate; this repository implements the
elements of the assignment and its immediate design context.

## Model mapping

| Model element | Type | Implementation |
| --- | --- | --- |
| DVBD.com | ApplicationComponent | `cmd/dvbd` (entrypoint) + `internal/webshop` |
| Webshop | ApplicationService | `internal/webshop` (Echo routes) |
| Webshop | BusinessService | `internal/webshop` pages (`/`, `/selectie`) |
| Product | DataObject | `internal/product` (`Product`, `Repository`) |
| Product | BusinessObject | `internal/product.Category` specializations |
| Fruit & Groenten | BusinessObject | `internal/product.CategoryFruitGroenten` |
| Pasta | BusinessObject | `internal/product.CategoryPasta` |
| Selecteer producten | BusinessProcess | `internal/selection` + `/selecteren`, `/api/selection` |
| Klant | BusinessActor | `internal/klant` (customer cookie) |
| Go/Echo | SystemSoftware | Go module + `github.com/labstack/echo/v4` |

The Product DataObject realises the Product BusinessObject, and DVBD.com
accesses the Product DataObject through the repository.

## Endpoints

Pages:

- `GET /` — the Webshop catalog, grouped by Fruit & Groenten and Pasta
- `GET /selectie` — the selection of the Klant
- `POST /selecteren` — select a product (Selecteer producten)
- `POST /selectie/verwijderen` — remove a product from the selection

JSON API:

- `GET /api/health`
- `GET /api/products`, `GET /api/products/:id`
- `GET /api/selection`, `POST /api/selection`, `DELETE /api/selection/:id`

## Persistence

The Product DataObject is stored in Postgres when a `DATABASE_URL` (or
`POSTGRES_DSN`) is configured. When no database is reachable, the Webshop falls
back to an in-memory catalog that is seeded from `internal/product/catalog.go`,
so the delivery can always be previewed. The local Postgres service is declared
in `.plot/package.json`.

## Development

```sh
npm run start   # start the Webshop on 0.0.0.0:8080
npm test        # run the Go tests
npm run build   # build the dvbd binary
```
