# 開発進行状況 (2025年10月22日)

## 現在の問題点

`openfga` サービスのヘルスチェックが `unhealthy` になり、`backend` サービスが起動しない。

## 原因の調査

*   `docker-compose ps` で `openfga` が `unhealthy` であることを確認。
*   `openfga` のログ (`docker-compose logs openfga`) で、gRPC サーバーがコンテナ内で設定されたアドレスでリッスンしていることを確認。

## 試したアプローチと結果

1.  **`openfga` のヘルスチェックコマンドを `wget` から `curl` に変更。**
    *   結果: `openfga` イメージに `curl` が含まれていないためビルド失敗。
2.  **`openfga.Dockerfile` を作成し `curl` をインストールしてビルド。**
    *   結果: `openfga` イメージが `alpine` ベースではないためビルド失敗 (`/bin/sh` がない)。
3.  **`grpc-health-probe` を使用するようにヘルスチェックコマンドを修正 (`-addr=localhost:8081`)。**
    *   結果: `unhealthy` のまま。`localhost` が IPv4 を指すため、IPv6 でリッスンしている `openfga` に接続できない可能性。
4.  **`grpc-health-probe` のアドレスを IPv4 ループバックアドレスに変更。**
    *   結果: `unhealthy` のまま。
5.  **`grpc-health-probe` のアドレスを IPv6 ループバックアドレスに変更。**
    *   結果: `unhealthy` のまま。Docker ネットワーク内で IPv6 ループバックへのルーティングがうまくいっていない可能性。
6.  **`openfga` の `command` オプションで gRPC アドレスを明示的に設定し、`healthcheck` を `grpc-health-probe` でコンテナ内のループバックアドレスをチェックするように戻す。**
    *   結果: `unhealthy` のまま。
7.  **`openfga` の `healthcheck` の `start_period` を `30s` に増やす。**
    *   結果: `unhealthy` のまま。
8.  **`openfga` コンテナ内で `grpc-health-probe` を手動実行。**
    *   結果: `status: SERVING` で正常終了。`openfga` サービス自体は正常に起動しており、`grpc-health-probe` も動作していることを確認。
9.  **`openfga` の `healthcheck` の `test` コマンドを `["CMD", "/usr/local/bin/grpc_health_probe", "-addr=[::]:8081"]` に修正。**
    *   結果: `openfga` サービスは `Healthy` になったが、`backend` サービスが `./server` が見つからないエラーで失敗。

## 現在の主要な問題点

`backend` サービスが `exec: "/app/server": stat /app/server: no such file or directory: unknown` エラーで起動しない。

## 原因の分析

*   `backend/Dockerfile` で実行ファイルを `/usr/local/bin/server` にビルドしている。
*   `docker-compose.yml` の `backend` サービスで `volumes: - ./backend:/app` が設定されており、ホストの `./backend` がコンテナの `/app` にマウントされる。
*   `Dockerfile` の `CMD ["./server"]` は `WORKDIR /app` の相対パスで `./server` を探すが、`/app` は `volumes` で上書きされているため、ビルドされた実行ファイルが見つからない。

## 解決策

1.  **`docker-compose.yml` の `backend` サービスから `volumes: - ./backend:/app` の行を削除。**
    *   `backend` サービスは `Dockerfile` でビルドされた実行ファイルを使用するため、ホストのソースコードをマウントする必要がない。
2.  **`backend/Dockerfile` の `CMD ["./server"]` を `CMD ["/app/server"]` に修正し、実行ファイルの絶対パスを指定。**
    *   実行ステージの `WORKDIR` が `/app` であり、ビルドされた実行ファイルが `/app/server` にコピーされるため、`CMD` もそれに合わせる。
3.  **`backend/Dockerfile` に `grpc-health-probe` をインストールする行を追加し、実行ステージにコピー。**
    *   `backend` サービスのヘルスチェックに `grpc-health-probe` を使用するため、コンテナ内に `grpc-health-probe` が必要。
4.  **`backend/main.go` に gRPC ヘルスチェックサービスを実装。**
    *   `backend` サービスが `grpc-health-probe` からのヘルスチェックに応答できるように、gRPC ヘルスチェックサービスを実装する。
5.  **`docker-compose.yml` の `backend` サービスに `healthcheck` を追加し、`grpc-health-probe` でコンテナ内のアドレスをチェックするように設定。**
    *   `backend` サービスが IPv6 のアドレスでリッスンしているため、ヘルスチェックもそれに合わせる。
6.  **`docker-compose.yml` の `backend` サービスの `healthcheck` の `start_period` を `10s` に設定。**

## 結果

すべてのサービス (`postgres`, `openfga`, `backend`) が `healthy` 状態で起動しました。

### ネットワーク設定とプロセスに関する補足

*   **`openfga` サービス:** gRPCサーバーはコンテナ内で設定されたアドレスでリッスンしています。ヘルスチェックには `grpc-health-probe` を使用し、`openfga` コンテナ内で `grpc-health-probe` が正常に動作することを確認しました。
*   **`backend` サービス:** gRPCサーバーはコンテナ内で設定されたアドレスでリッスンしています。ヘルスチェックには `grpc-health-probe` を使用し、`backend` コンテナ内に `grpc-health-probe` をインストールし、`backend/main.go` にヘルスチェックサービスを実装することで、正常に動作するようになりました。