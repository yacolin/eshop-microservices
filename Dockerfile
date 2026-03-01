# 多阶段构建，支持多服务：通过 build arg SERVICE 指定 cmd 下的服务名
# 仅改 .go 代码时：go mod download 层会命中缓存，不会重新拉依赖；需启用 BuildKit：DOCKER_BUILDKIT=1
# 例：docker build --build-arg SERVICE=order-service -t order-service .
ARG SERVICE=order-service
FROM golang:1.21-alpine AS builder
ARG SERVICE
WORKDIR /app

# 先只复制依赖文件，此层不变则 go mod download 使用缓存
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-bin ./cmd/${SERVICE}

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app-bin ./app
COPY configs ./configs
EXPOSE 8080
CMD ["./app"]
