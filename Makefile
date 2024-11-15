postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=shu -e POSTGRES_PASSWORD=shu -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=shu --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=shu --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://shu:shu@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://shu:shu@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb sqlc test migrateup migratedown server