postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root go_transact
dropdb:
	docker exec -it postgres16 dropdb go_transact
migrateup:
	migrate -path db/migration  -database "postgresql://root:password@localhost:5432/go_transact?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration  -database "postgresql://root:password@localhost:5432/go_transact?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration  -database "postgresql://root:password@localhost:5432/go_transact?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration  -database "postgresql://root:password@localhost:5432/go_transact?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/petershivachi/go_transact/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1