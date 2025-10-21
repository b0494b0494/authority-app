# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Goモジュールファイルをコピーして依存関係をダウンロード
COPY backend/go.mod .
COPY backend/go.sum .
RUN go mod download

# Install protobuf compiler and Go plugin
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# protoファイルをコピー
COPY proto ./proto

# Generate Go code from .proto files
RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/directory/v1/directory.proto

# ソースコードをコピー
COPY backend/main.go .

# アプリケーションをビルド
RUN go build -o /usr/local/bin/server .

# Install grpc-health-probe
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# ビルドステージから実行可能ファイルをコピー
COPY --from=builder /usr/local/bin/server .
COPY --from=builder /go/bin/grpc-health-probe /usr/local/bin/grpc_health_probe

# 必要な証明書をインストール (gRPCのTLS/SSL接続用など)
RUN apk add --no-cache ca-certificates

# アプリケーションを実行
CMD ["/app/server"]