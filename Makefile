postgres:
	docker run --name some-postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it some-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it some-postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:aolin940118@localhost:5432/simple_bank?sslmode=disable" --verbose up
migrateup1:
	migrate -path db/migration -database "postgres://root:aolin940118@localhost:5432/simple_bank?sslmode=disable" --verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://root:aolin940118@localhost:5432/simple_bank?sslmode=disable" --verbose down
migratedown1:
	migrate -path db/migration -database "postgres://root:aolin940118@localhost:5432/simple_bank?sslmode=disable" --verbose down 1

migrateforce:

	migrate -path db/migration -database "postgres://root:aolin940118@localhost:5432/simple_bank?sslmode=disable" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...
server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/aolin118/simple_bank/db/sqlc Store 
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock