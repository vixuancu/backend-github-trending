# Dockerfile cho ứng dụng backend-github-trending
# Sử dụng multi-stage build để giảm kích thước image cuối cùng

# Stage 1: Build stage - Giai đoạn xây dựng
# Sử dụng phiên bản Go 1.23.0 với Alpine Linux làm base image
FROM golang:1.23.0-alpine AS builder

# Cài đặt các dependency cần thiết để build
# git: để go có thể tải các package từ các repository git
RUN apk update && apk add --no-cache git

# Thiết lập thư mục làm việc trong container
WORKDIR /app

# Copy file go.mod và go.sum trước để tận dụng Docker cache
# Các layer này sẽ không bị rebuild nếu các file này không thay đổi
COPY go.mod go.sum ./

# Tải các dependency
RUN go mod download

# Copy toàn bộ source code vào container
COPY . .

# Build ứng dụng
# CGO_ENABLED=0: tắt CGO để build binary hoàn toàn static
# GOOS=linux: build cho hệ điều hành Linux
# -a: force rebuild tất cả package
# -installsuffix cgo: thêm hậu tố vào package installation directory
# -o backend-github-trending: tên file output
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend-github-trending .

# Stage 2: Final stage - Giai đoạn cuối cùng
# Sử dụng Alpine Linux nhẹ nhất có thể làm base image cho container cuối cùng
FROM alpine:latest

# Cài đặt các package cần thiết cho runtime
# ca-certificates: để hỗ trợ HTTPS
# tzdata: để hỗ trợ múi giờ
RUN apk --no-cache add ca-certificates tzdata

# Thiết lập thư mục làm việc
WORKDIR /app

# Copy file binary đã build từ stage 1
COPY --from=builder /app/backend-github-trending .

# Copy các file cấu hình cần thiết
COPY --from=builder /app/dbconfig.yml .
# Nếu có thêm file cấu hình khác, bạn có thể copy thêm ở đây

# Tạo thư mục log nếu ứng dụng cần
#RUN mkdir -p /app/log_files/error /app/log_files/info

# Thiết lập các biến môi trường (nếu cần)


# Expose cổng mà ứng dụng sẽ chạy (giả sử ứng dụng chạy ở cổng 8080)
EXPOSE 8080

# Lệnh chạy ứng dụng khi container khởi động
CMD ["./backend-github-trending"]
