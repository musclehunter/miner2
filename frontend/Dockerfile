FROM node:lts-alpine

WORKDIR /app

# 依存関係のインストール
COPY package.json ./
# package-lock.jsonが存在しない場合でもインストールできるようにする
RUN npm install

# アプリケーションのコピー
COPY . .

EXPOSE 8080

# 開発サーバーの起動
CMD ["npm", "run", "serve"]
