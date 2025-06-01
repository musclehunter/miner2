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
│   simple_server.go  # エントリーポイント
│   handlers/         # ハンドラー関数
│   ├── auth.go       # 認証関連ハンドラー
│   ├── game.go       # ゲーム関連ハンドラー
│   ├── admin.go      # 管理者用ハンドラー
│   └── pending_users.go # 未確認ユーザー管理
│   models/           # データモデル
│   database/         # データベース関連
│   cache/            # Redisキャッシュ関連
└── mail/             # メール送信関連
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

## 機能概要

### メール認証システム

ユーザー登録時にメール認証を実行します。主なフローは以下のとおりです：

1. ユーザーがメールアドレスとパスワードで登録を実行
2. ユーザー情報がRedisに仮登録され、メール確認トークンが生成される
3. メール確認リンクが送信される
4. ユーザーがメール内のリンクをクリックしてメールアドレスを確認
5. 確認後、Redisから情報が取得され、正式にデータベースにユーザーが登録される

関連エンドポイント：
- `POST /api/auth/signup`: ユーザー登録と確認メール送信
- `GET /api/auth/verify-email`: メール確認トークンの検証
- `POST /api/auth/resend-verification`: 確認メールの再送信

### 管理者API

管理者向けのAPIエンドポイントが実装されています。すべての管理者APIは `AdminAuth` ミドルウェアによって保護されています。

主な管理者API機能：

#### ユーザー管理
- `GET /api/admin/users`: 全ユーザー一覧取得
- `GET /api/admin/users/:id`: ユーザー詳細取得
- `PUT /api/admin/users/:id`: ユーザー情報更新
- `DELETE /api/admin/users/:id`: ユーザー削除

#### 未確認ユーザー管理
- `GET /api/admin/pending-users`: メール未確認ユーザー一覧取得

#### 町データ管理
- `GET /api/admin/towns`: 町一覧取得
- `POST /api/admin/towns`: 新規町作成
- `PUT /api/admin/towns/:id`: 町情報更新
- `DELETE /api/admin/towns/:id`: 町削除

### Redisキャッシュ

Redisは主に以下の目的で使用されています：

1. メール確認トークンの保存
2. 未確認ユーザーの仮登録情報の管理

開発環境では、環境変数 `USE_MOCK_REDIS` でRedisクライアントの使用を制御できます：

- `USE_MOCK_REDIS=true`: メモリ内モックRedisクライアントを使用
- `USE_MOCK_REDIS=false`: Docker環境内の実際Redisサーバーに接続

デフォルトでは `docker-compose.yml` で `USE_MOCK_REDIS=false` が設定されており、実際Redisサーバーを使用します。本番環境でも同様に実際Redisサーバーに接続します。

Redis接続には以下の環境変数を使用します：

- `REDIS_HOST`: Redisサーバーのホスト名（デフォルト: `redis`）
- `REDIS_PORT`: Redisサーバーのポート（デフォルト: `6379`）
- `REDIS_PASSWORD`: Redisのパスワード（デフォルト: 空）
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
