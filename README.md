## How to run ?
Run with 
```sh
$ GOCHAT_ENV=development go run ./cmd/gochat
```
to pretty print logs during development

To create a migration
```
migrate create -ext sql -dir internal/database/migrations -seq <name>
```

To run migrations
```
migrate -database ${GOCHAT_DB_DSN} -path internal/database/migrations up
```

To run using docker,
```
docker run --env-file .env --network host ananthvk0/gochat:0.0.1
```

```
$ docker volume create pg_data
$ docker run -e POSTGRES_PASSWORD=dev -p 0.0.0.0:5432:5432 -v pg_data:/var/lib/postgresql/data postgres:16-alpine
```


## TODO
- [ ] Clients receive messages from all rooms in the same chat, fix that