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
	migrate -path ./migrations/postgres -database 'postgres://abdulahad:passwd123@localhost:5432/article_service_db?sslmode=disable' up

migratedown:
	migrate -path ./migrations/postgres -database 'postgres://abdulahad:passwd123@localhost:5432/article_service_db?sslmode=disable' down

proto-author:
	protoc --go_out=./genproto --go-grpc_out=./genproto protos/author_service/*.proto
proto-article:
	protoc --go_out=./genproto --go-grpc_out=./genproto protos/article_service/*.proto




