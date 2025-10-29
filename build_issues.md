 # ビルドの問題点

## ローカルビルド

- `go build` を実行すると、`openfga/go-sdk` のAPIの使い方が古いためにコンパイルエラーが発生しました。
- `openfga_client.go` と `upload_service.go` を修正し、`go-sdk v0.7.3` のAPIに合わせる必要がありました。
- `go-sqlite3` を使っているため、`CGO_ENABLED=1` でビルド・実行する必要がありましたが、ローカル環境の `go env` で `CGO_ENABLED=0` となっていたため、実行時にエラーが発生しました。

## Dockerビルド

- `Dockerfile` で `CGO_ENABLED=1` を設定していましたが、`golang:1.24-alpine` イメージに `gcc` が含まれていないため、ビルドに失敗しました。
- `gcc` をインストールすると、`stdlib.h` が見つからないというエラーが発生しました。
- `musl-dev` をインストールすることで、Cの標準ライブラリの問題は解決しました。
- `openfga/go-sdk` のAPIの使い方が古いために、`openfga_client.go` でコンパイルエラーが発生しました。
- `Userset` の `This` フィールドが存在しない、`Users` フィールドが存在しない、などの問題がありました。
