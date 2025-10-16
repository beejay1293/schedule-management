MIGRATION_PATH=backend/internal/db/migrations

.PHONY: create-migration gen-proto start-app

# Create a new migration file
create-migration:
	@if [ -z "${name}" ]; then \
		echo "Error: please provide a migration name (e.g. make create-migration name=add_appointments_table)"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}

# Generate Go protobuf & gRPC files
gen-proto:
	protoc \
		-I=proto \
		--go_out=internal/pb --go_opt=paths=source_relative \
		--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
		proto/appointment.proto

# Start the Go backend app
start-app:
	go run cmd/server/main.go
