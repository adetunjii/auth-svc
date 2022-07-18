run: |
	gofmt -w .
	go run main.go

mock-service-user:
	mockgen -source=service/user.go -destination=service/user_mock.go -package=service

gen:
	protoc --go_out=internal/proto --go_opt=paths=source_relative \
    --go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative \
    --proto_path=internal/proto internal/proto/*.proto
