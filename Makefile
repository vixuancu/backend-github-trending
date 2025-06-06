pro:
	docker rmi -f backend-github-trending:latest
	docker-compose up

dev:
	cd cmd/dev; go run main.go

# Chạy ứng dụng từ file main.go chính
run:
	go run main.go

# Chạy Docker Compose ở chế độ detach (chạy ngầm)
prod-d:
	docker rmi -f backend-github-trending:latest
	docker-compose up -d

# Dừng và xóa các container
down:
	docker-compose down

# Xây dựng lại Docker image và chạy
rebuild:
	docker-compose down
	docker rmi -f backend-github-trending:latest
	docker-compose up --build

# Thêm lệnh để dọn dẹp images không sử dụng
clean-docker:
	docker system prune -f

# Xóa các image cũ và dangling
clean-images:
	docker image prune -a -f

# Build ứng dụng
build:
	go build -o bin/github-trending main.go

# Chạy tất cả các test
test:
	go test -v ./...

# Dọn dẹp các file build
clean:
	rm -rf bin/
	go clean

# Hiển thị thông tin trợ giúp
help:
	@echo "Các lệnh Makefile:"
	@echo "  pro        : Khởi động lại ứng dụng với Docker Compose"
	@echo "  prod-d     : Khởi động lại ứng dụng với Docker Compose ở chế độ detach"
	@echo "  dev        : Chạy ứng dụng ở chế độ dev"
	@echo "  run        : Chạy ứng dụng từ main.go"
	@echo "  build      : Build ứng dụng thành file thực thi"
	@echo "  down       : Dừng và xóa các container"
	@echo "  rebuild    : Xây dựng lại Docker image từ đầu và chạy"
	@echo "  clean-docker: Dọn dẹp các Docker images không sử dụng"
	@echo "  test       : Chạy tất cả các test"
	@echo "  clean      : Xóa thư mục bin và dọn dẹp"
