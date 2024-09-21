up:
	docker compose up -d

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop

in:
	docker compose exec app bash

indb:
	docker compose exec db bash

init:
	go mod tidy
	go build -o gent cmd/gent/main.go

run:
	go run cmd/shiftmanager/main.go

test:
	go test
