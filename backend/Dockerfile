FROM golang:1.21-alpine

WORKDIR /app

# 依存関係のインストール
COPY go.mod ./
# go.sumはまだ存在しないかもしれないため、後で作成されるようにする
RUN go mod download && go mod tidy

# ホットリロード用のcompileDaemon（開発環境用）
RUN go install github.com/githubnemo/CompileDaemon@latest

# アプリケーションのコピー
COPY . .

EXPOSE 8080

# CompileDaemonを使ってホットリロード
CMD CompileDaemon --build="go build -o ./tmp/main ." --command="./tmp/main"
