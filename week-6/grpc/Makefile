.PHONY: gen-proto
gen-proto:
	@protoc \
		--proto_path=protobuf protobuf/*.proto \
		--go_out=services/common/genproto/orders --go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/orders --go-grpc_opt=paths=source_relative

.PHONY: run-orders
run-orders:
	@go run services/orders/cmd/main.go

.PHONY: run-kitchen
run-kitchen:
	@go run services/kitchen/cmd/main.go
