# 阶段1：用官方 Go 镜像编译
FROM golang:1.25-alpine AS builder

WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./

# 国内 Go 代理，否则 proxy.golang.org 连不上
ENV GOPROXY=https://goproxy.cn,direct

RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 go build -o server main.go

# 阶段2：用极简镜像运行
FROM alpine:latest

WORKDIR /app

# 安装 CA 证书（HTTPS 请求需要）
RUN apk --no-cache add ca-certificates tzdata

# 从编译阶段复制二进制和配置文件
COPY --from=builder /build/server .
COPY --from=builder /build/.env .

EXPOSE 8080

CMD ["./server"]
