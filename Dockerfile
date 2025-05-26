# ใช้ base image ของ Go สำหรับ multi-platform
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# ตั้ง working directory
WORKDIR /app

# ติดตั้ง dependencies ที่จำเป็น
RUN apk add --no-cache git

# คัดลอก go.mod และ go.sum สำหรับ cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# คัดลอกโค้ดทั้งหมด
COPY . .

# Build สำหรับ ARM64 (M1/M2)
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .


# ใช้ base image ที่เล็กสำหรับรันแอพ
FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /app

# ติดตั้งแพคเกจที่จำเป็น
RUN apk add --no-cache ca-certificates tzdata

# คัดลอก binary และ config
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# ตั้งค่า timezone ไทย
ENV TZ=Asia/Bangkok

# ตั้งค่าสิทธิ์การรัน
RUN chmod +x /app/main

# กำหนด port
EXPOSE 8080

# คำสั่งรันแอพ
CMD ["./main"]