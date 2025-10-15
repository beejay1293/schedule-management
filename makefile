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

gen-frontend-proto:
	protoc \
		--js_out=import_style=commonjs,binary:./src/pb \
		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:./src/pb \
		appointment.proto

run-backend:
	cd backend && go run cmd/server/main.go

run-http:
	cd backend && go run http_server.go

run-grpc-proxy:
	grpcwebproxy \
		--backend_addr=localhost:50051 \
		--run_tls_server=false \
		--allow_all_origins \
		--server_http_debug_port=8080
