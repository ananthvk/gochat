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
migrate -database ${POSTGRESQL_URL} -path internal/database/migrations up
```

To run using docker,
```
docker run --env-file .env --network host ananthvk0/gochat:0.0.1
```


## TODO
- [ ] Clients receive messages from all rooms in the same chat, fix that