.PHONY: help build run test clean docker-build docker-run k8s-deploy k8s-delete

help:
	@echo "API Server Makefile Commands:"
	@echo "  make build         - 编译 Go 应用"
	@echo "  make run           - 运行应用"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建文件"
	@echo "  make docker-build  - 构建 Docker 镜像"
	@echo "  make docker-run    - 使用 Docker Compose 运行"
	@echo "  make k8s-deploy    - 部署到 Kubernetes"
	@echo "  make k8s-delete    - 从 Kubernetes 删除"
	@echo "  make deps          - 下载依赖"

build:
	@echo "Building API Server..."
	go build -o bin/api-server ./cmd/api-server

run:
	@echo "Running API Server..."
	go run ./cmd/api-server/main.go

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

docker-build:
	@echo "Building Docker image..."
	docker build -t api-server:latest .

docker-run:
	@echo "Running with Docker Compose..."
	docker-compose -f deployments/docker/docker-compose.yml up -d

docker-stop:
	@echo "Stopping Docker Compose..."
	docker-compose -f deployments/docker/docker-compose.yml down

k8s-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -k deployments/kubernetes/

k8s-delete:
	@echo "Deleting from Kubernetes..."
	kubectl delete -k deployments/kubernetes/

k8s-logs:
	@echo "Showing logs..."
	kubectl logs -n api-server -l app=api-server --tail=100 -f

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

dev:
	@echo "Running in development mode with auto-reload..."
	air
