FROM golang:1.21-alpine

WORKDIR /app

# 必要なシステムパッケージをインストール
RUN apk add --no-cache gcc musl-dev git

# プロジェクトファイル全体をコピー
COPY . .

# 依存関係を直接インストール
RUN go mod tidy \
    && go get -u github.com/go-sql-driver/mysql@v1.7.1 \
    && go get -u github.com/google/uuid@v1.3.0 \
    && go get -u github.com/golang-jwt/jwt/v5@v5.0.0 \
    && go mod tidy

# ホットリロード用のcompileDaemon（開発環境用）
RUN go install github.com/githubnemo/CompileDaemon@latest

# アプリケーションのコピー
COPY . .

EXPOSE 8080

# CompileDaemonを使ってホットリロード
CMD CompileDaemon --build="go build -o ./tmp/main ." --command="./tmp/main"
