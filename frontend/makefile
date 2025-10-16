.PHONY: protogen
protogen: ./proto/
	docker build --tag protogen:latest .
	docker run -v $(CURDIR)/proto:/proto:ro -v $(CURDIR)/src/grpc-web-client-gen:/gen protogen:latest

gen-frontend-proto:
	protoc \
		--proto_path=./backend/proto \
		--js_out=import_style=commonjs,binary:./frontend/src/pb \
		--grpc-web_out=import_style=typescript,mode=grpcweb:./frontend/src/pb \
		./backend/proto/appointment.proto

run-grpc-proxy:
	grpcwebproxy \
		--backend_addr=localhost:50051 \
		--run_tls_server=false \
		--allow_all_origins \
		--server_http_debug_port=8080
