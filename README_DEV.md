# 採掘人 - 開発ガイド

## 開発環境セットアップ

このプロジェクトはdocker-composeを使用して開発環境を構築しています。

### 前提条件

- Docker
- Docker Compose

### 環境構築

1. リポジトリをクローン（または既存のフォルダで実行）:

```bash
git clone <repository-url>
cd miner2
```

2. Docker Composeでサービスを起動:

```bash
docker-compose up -d
```

これにより以下のサービスが起動します:
- バックエンド (Go/Gin): http://localhost:8080
- フロントエンド (Vue.js): http://localhost:8081
- MySQL: localhost:3306
- Redis: localhost:6379

### バックエンド開発 (Go/Gin)

バックエンドコードは `/backend` ディレクトリにあります。
ホットリロード対応のため、コードを変更すると自動的に再起動します。

#### 主要なディレクトリ構造:

```
backend/
├── main.go           # エントリーポイント
├── routes/           # ルート定義
├── controllers/      # コントローラー
├── models/           # データモデル
├── middleware/       # ミドルウェア
└── utils/            # ユーティリティ
```

### フロントエンド開発 (Vue.js)

フロントエンドコードは `/frontend` ディレクトリにあります。
ホットリロード対応のため、コードを変更すると自動的に再ビルドされます。

#### 主要なディレクトリ構造:

```
frontend/
├── public/           # 静的アセット
└── src/              # ソースコード
    ├── assets/       # アセット
    ├── components/   # コンポーネント
    ├── views/        # ビュー
    ├── router/       # ルーター
    ├── store/        # Vuex ストア
    └── App.vue       # ルートコンポーネント
```

### データベース操作

MySQLへの接続:

```bash
docker-compose exec db mysql -u miner -p minerdb
# パスワード: minerpassword
```

## コンテナの管理

すべてのサービスを起動:
```bash
docker-compose up -d
```

ログの確認:
```bash
docker-compose logs -f
```

特定のサービスのログ確認（例: バックエンド）:
```bash
docker-compose logs -f backend
```

すべてのサービスを停止:
```bash
docker-compose down
```

コンテナを再ビルド（依存関係の変更時）:
```bash
docker-compose build
```

## テスト

バックエンドのテスト実行:
```bash
docker-compose run --rm backend go test ./...
```

フロントエンドのテスト実行:
```bash
docker-compose run --rm frontend npm run test
```
