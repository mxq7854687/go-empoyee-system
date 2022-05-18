postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root employee

dropdb:
	docker exec -it postgres12 dropdb  employee

migrateup:
	migrate -path db/migration/ -database "postgresql://root:admin@localhost:5432/employee?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/ -database "postgresql://root:admin@localhost:5432/employee?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdxb dropdb migrateup migratedown sqlc test server