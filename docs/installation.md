# 採掘人 - インストールガイド

## 前提条件

このプロジェクトを実行するには以下のソフトウェアが必要です：

- Docker
- Docker Compose

## インストール手順

### 1. リポジトリのクローン

まず、GitHubリポジトリをクローンします：

```bash
git clone <repository-url>
cd miner2
```

### 2. 環境変数の設定

必要に応じて環境変数を設定します。デフォルト設定ではローカル環境での開発に最適化されています。

#### バックエンド環境変数

`.env`ファイルが必要な場合は、サンプルをコピーして作成します：

```bash
cp backend/.env.example backend/.env
```

主な環境変数：
- `DB_HOST` - MySQLホスト（デフォルト: db）
- `DB_PORT` - MySQLポート（デフォルト: 3306）
- `DB_USER` - MySQLユーザー（デフォルト: miner）
- `DB_PASSWORD` - MySQLパスワード（デフォルト: minerpassword）
- `DB_NAME` - データベース名（デフォルト: minerdb）
- `REDIS_HOST` - Redisホスト（デフォルト: redis）
- `REDIS_PORT` - Redisポート（デフォルト: 6379）
- `REDIS_PASSWORD` - Redisパスワード（デフォルト: 空）
- `USE_MOCK_REDIS` - Redisモックモード（true/false）

#### フロントエンド環境変数

必要に応じて、フロントエンドの環境変数を設定します：

```bash
cp frontend/.env.example frontend/.env
```

主な環境変数：
- `VUE_APP_API_URL` - バックエンドAPIのURL
- `VUE_APP_USE_MOCK` - モックモードの有効/無効（false推奨）

### 3. Dockerコンテナの起動

Docker Composeを使用してすべてのサービスを起動します：

```bash
docker-compose up -d
```

これにより以下のサービスが起動します：
- バックエンド (Go/Gin): http://localhost:8080
- フロントエンド (Vue.js): http://localhost:8081
- MySQL: localhost:3306
- Redis: localhost:6379

### 4. データベースのセットアップ

初回起動時にはデータベースの初期化が自動的に行われます。追加のセットアップが必要な場合：

```bash
docker-compose exec backend go run ./scripts/seed.go
```

### 5. アクセス方法

セットアップが完了したら、以下のURLでアプリケーションにアクセスできます：

- フロントエンド: http://localhost:8081
- バックエンドAPI: http://localhost:8080

## トラブルシューティング

### コンテナが起動しない

ログを確認して問題を特定します：

```bash
docker-compose logs
```

特定のサービスのログを確認：

```bash
docker-compose logs backend
docker-compose logs frontend
```

### データベース接続エラー

MySQLサービスが起動しているか確認：

```bash
docker-compose ps
```

MySQLコンテナに接続してデータベースを確認：

```bash
docker-compose exec db mysql -u miner -p minerdb
# パスワード: minerpassword
```

### Redis接続エラー

Redisサービスが起動しているか確認：

```bash
docker-compose exec redis redis-cli ping
```

レスポンスが「PONG」であればRedisは正常に動作しています。
