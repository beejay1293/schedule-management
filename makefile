.PHONY: gen-proto run-backend run-http

gen-proto:
	cd backend && protoc --go_out=. --go-grpc_out=. proto/appointment.proto

run-backend:
	cd backend && go run main.go

run-http:
	cd backend && go run http_server.go
