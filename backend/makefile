MIGRATION_PATH=backend/internal/db/migrations
DB_URL=postgres://postgres:password@localhost:5432/mydb?sslmode=disable

create-migration:
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}

migrate-up:
	migrate -path ${MIGRATION_PATH} -database ${DB_URL} up

migrate-down:
	migrate -path ${MIGRATION_PATH} -database ${DB_URL} down

.PHONY: gen-proto run-backend run-http

gen-backend-proto:
	protoc \
		-I=backend/proto \
		--go_out=backend/internal/pb --go_opt=paths=source_relative \
		--go-grpc_out=backend/internal/pb --go-grpc_opt=paths=source_relative \
		appointment.proto

run-backend:
	go run cmd/server/main.go

