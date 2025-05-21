# idm
OOP in Go

## Testing

go test ./...

go test -cover ./...

golangci-lint run

## Migrations

goose create <script_name> sql

goose up -dir ./migrations

## Swagger

swag init -d cmd,inner --parseDependency --parseInternal

https://localhost:8080/swagger

## Authentication

https://www.keycloak.org/app/#url=http://localhost:9990&realm=idm&client=idmapp

admin: new-user@idm.ru / 12345
