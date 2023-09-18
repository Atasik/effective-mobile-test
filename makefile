.PHONY:
.SILENT:

build:
    go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/fio ./cmd/fio/main.go 

run: build
    docker compose up fio

test:
    go test -v ./...

migrate:
    migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/fioDb?sslmode=disable' up

kafka-ui:
    docke compose up kafka-ui

swag:
    swag init -g internal/app/app.go