
# Go POC

Proof of Concept

* [Hexagonal Architecture](concept/hexagonal/README.md)
* [Batch Fetching](concept/batchfetching/README.md)
* [Database Transaction](concept/dbtransaction/README.md)

## Run App
```
$ go run main.go
```

## Create Environment
```
$ cp .env-example .env
```

## Start Docker Compose Deployment
```
$ docker-compose -f docker-compose.yml up -d --build
```

## Stop Docker Compose Deployment
```
$ docker-compose -f docker-compose.yml down
```

## Rebuild Docker Compose App Only
```
$ docker-compose up -d --no-deps --build poc
```

### Create Migration
```
$ migrate create -ext sql -dir service/{service_name}/migration/ -seq init_mg
```

### Load Test
```
$ k6 run -e MY_HOSTNAME=http://localhost:8000 loadtest/name.js
```