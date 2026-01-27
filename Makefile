.PHONY: run build clean test fmt lint help

# 默认目标
.DEFAULT_GOAL := help

# 运行项目
run:
	@echo "Starting CharonOMS..."
	@go run cmd/server/main.go

# 编译项目
build:
	@echo "Building CharonOMS..."
	@CGO_ENABLED=0 go build -o bin/charonoms cmd/server/main.go
	@echo "Build completed: bin/charonoms"

# Linux 编译
build-linux:
	@echo "Building CharonOMS for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/charonoms-linux cmd/server/main.go
	@echo "Build completed: bin/charonoms-linux"

# Windows 编译
build-windows:
	@echo "Building CharonOMS for Windows..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/charonoms.exe cmd/server/main.go
	@echo "Build completed: bin/charonoms.exe"

# 清理编译文件
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f charonoms charonoms.exe
	@echo "Clean completed"

# 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format completed"

# 代码检查
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Lint completed"

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Dependencies downloaded"

# 更新依赖
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"

# 生成 go.sum
tidy:
	@echo "Tidying go.mod..."
	@go mod tidy
	@echo "Tidy completed"

# Docker 构建
docker-build:
	@echo "Building Docker image..."
	@docker build -t charonoms:latest .
	@echo "Docker build completed"

# Docker 运行
docker-run:
	@echo "Running Docker container..."
	@docker run -d -p 5001:5001 --name charonoms charonoms:latest
	@echo "Docker container started"

# Docker 停止
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop charonoms
	@docker rm charonoms
	@echo "Docker container stopped"

# 帮助信息
help:
	@echo "CharonOMS Makefile Commands:"
	@echo ""
	@echo "  make run            - 运行项目"
	@echo "  make build          - 编译项目"
	@echo "  make build-linux    - 编译 Linux 版本"
	@echo "  make build-windows  - 编译 Windows 版本"
	@echo "  make clean          - 清理编译文件"
	@echo "  make test           - 运行测试"
	@echo "  make fmt            - 格式化代码"
	@echo "  make lint           - 代码检查"
	@echo "  make deps           - 下载依赖"
	@echo "  make deps-update    - 更新依赖"
	@echo "  make tidy           - 整理 go.mod"
	@echo "  make docker-build   - 构建 Docker 镜像"
	@echo "  make docker-run     - 运行 Docker 容器"
	@echo "  make docker-stop    - 停止 Docker 容器"
	@echo "  make help           - 显示帮助信息"
