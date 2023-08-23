gen-doc:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g internal/api/handler.go 
	swag fmt

lint:
	go fmt ./...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --fix

build:
	docker build -t $(tag) -f Dockerfile .