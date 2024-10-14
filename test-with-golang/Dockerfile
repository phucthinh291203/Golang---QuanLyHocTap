# Sử dụng image Go 1.23
FROM golang:1.23-alpine AS builder

# Thiết lập thư mục làm việc trong container
WORKDIR /app

# Sao chép go.mod và go.sum vào thư mục làm việc
COPY go.mod go.sum ./

# Tải các phụ thuộc
RUN go mod download

# Sao chép mã nguồn vào thư mục làm việc
COPY . .

# Biên dịch ứng dụng
RUN go build -o myapp .

# Chạy ứng dụng
CMD ["./myapp"]