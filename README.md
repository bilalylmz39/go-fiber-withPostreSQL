# Products API

Katmanli Go Fiber + PostgreSQL urun API'si. Proje, DDD'ye yakin bir ayrimla `internal` altinda domain, use-case, infrastructure ve HTTP adapter katmanlarina bolundu.

## Mimari

```txt
internal/
  domain/product          Entity, domain errors, repository contract
  application/product     Business use-cases and unit tests
  infrastructure/postgres PostgreSQL adapter
  interfaces/httpapi      Fiber routes and DTO parsing
  config                  Environment based configuration
```

`main.go` sadece dependency wiring yapar; business logic HTTP veya SQL paketlerini import etmez.

## Calistirma

PostgreSQL'de tabloyu olusturun:

```sh
psql -d productsdb -f schema.sql
```

Uygulamayi baslatin:

```sh
go run .
```

Varsayilan ortam degiskenleri:

```txt
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123
DB_NAME=productsdb
DB_SSLMODE=disable
PORT=3000
REQUEST_TIMEOUT=5
```

Alternatif olarak tek parca `DATABASE_URL` verebilirsiniz.

## Endpoints

```txt
GET    /health
GET    /products?query=book&active=true&min_price=10&max_price=100&limit=20&offset=0
GET    /products/:id
POST   /products
PUT    /products/:id
PATCH  /products/:id/stock
PATCH  /products/:id/activate
PATCH  /products/:id/deactivate
DELETE /products/:id
```

Urun payload'u:

```json
{
  "sku": "BK-001",
  "title": "Golang Book",
  "description": "It's a good book",
  "price": 42.12,
  "stock": 10
}
```

Stok guncelleme:

```json
{
  "delta": -2
}
```

## Test

```sh
go test ./...
```
