# 第一阶段：构建
FROM golang:1.23.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# 第二阶段：运行
FROM alpine

WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/config.dev.yaml .

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENTRYPOINT ["./app", "server"]