swag-init:
	swag init -g api/api.go -o api/docs
	
run: 
	go run cmd/main.go

install:
	swag init -g api/api.go -o api/docs
	go mod download
	go mod vendor
	go run cmd/main.go

migrateup:
	migrate -path ./migrations/postgres -database 'postgres://abdulahad:passwd123@localhost:5432/uacademy?sslmode=disable' up

migratedown:
	migrate -path ./migrations/postgres -database 'postgres://abdulahad:passwd123@localhost:5432/uacademy?sslmode=disable' down

proto:
	protoc --go_out=./genproto --go-grpc_out=./genproto protos/author_service/*.proto
