# 採掘人 - 開発ガイド

## 開発環境

このプロジェクトは開発環境として Docker Compose を使用し、ホットリロードに対応したコンテナ構成になっています。

## コンテナ管理

### 起動と停止

すべてのサービスを起動：
```bash
docker-compose up -d
```

すべてのサービスを停止：
```bash
docker-compose down
```

ログの確認：
```bash
docker-compose logs -f
```

特定のサービスのログ確認（例: バックエンド）：
```bash
docker-compose logs -f backend
```

コンテナを再ビルド（コードの大きな変更後など）：
```bash
docker-compose build
```

## バックエンド開発 (Go/Gin)

### ディレクトリ構造

バックエンドコードは `/backend` ディレクトリにあります：

```
backend/
│   simple_server.go  # エントリーポイント
│   handlers/         # APIハンドラー
│   models/           # データモデル
│   database/         # データベース操作
│   middleware/       # ミドルウェア
│   cache/            # Redisキャッシュ
└── mail/             # メール機能
```

### ホットリロード

バックエンドは CompileDaemon を使用したホットリロードに対応しています。コードを変更すると自動的に再コンパイルと再起動が行われます。

### テスト

バックエンドのテスト実行：
```bash
docker-compose run --rm backend go test ./...
```

特定のパッケージのテスト実行：
```bash
docker-compose run --rm backend go test ./handlers
```

### データベース操作

MySQLへの接続：
```bash
docker-compose exec db mysql -u miner -p minerdb
# パスワード: minerpassword
```

よく使うSQLコマンド：
```sql
-- テーブル一覧表示
SHOW TABLES;

-- ユーザーテーブルの内容確認
SELECT * FROM users;

-- 町のデータ確認
SELECT * FROM towns;
```

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

Redisデータの確認：
```bash
docker-compose exec redis redis-cli keys "*"
```

特定のキーの値を取得：
```bash
docker-compose exec redis redis-cli get "キー名"
```

## フロントエンド開発 (Vue.js)

### ディレクトリ構造

フロントエンドコードは `/frontend` ディレクトリにあります：

```
frontend/
├── public/           # 静的アセット
└── src/              # ソースコード
    ├── assets/       # 画像、スタイルなど
    ├── components/   # 再利用可能なコンポーネント
    ├── views/        # ページコンポーネント
    ├── router/       # Vue Router設定
    ├── store/        # Vuex状態管理
    ├── services/     # APIサービス
    └── App.vue       # ルートコンポーネント
```

### ホットリロード

フロントエンドもホットリロードに対応しています。コードを変更すると自動的に再コンパイルされ、ブラウザが更新されます。

### テスト

フロントエンドのテスト実行：
```bash
docker-compose run --rm frontend npm run test
```

### 認証システム

フロントエンドの認証システムは以下の機能を提供しています：

1. **ルートガード**
   - Vue Routerのナビゲーションガードにより、認証が必要なルートを保護
   - すべてのゲーム関連ルートに`requiresAuth: true`を設定

2. **認証状態管理**
   - Vuexストアで認証状態を管理
   - ローカルストレージにJWTトークンを保存

3. **API通信**
   - axios インスタンスで認証ヘッダーを自動設定
   - 401エラーでログアウト処理を実行

## APIリファレンス

### 認証関連

- `POST /api/auth/signup` - ユーザー登録
- `POST /api/auth/login` - ログイン
- `GET /api/auth/me` - 現在のユーザー情報取得
- `GET /api/auth/verify-email` - メールアドレス確認
- `POST /api/auth/resend-verification` - 確認メール再送信

### ゲーム関連

- `GET /api/game/towns` - 町の一覧取得
- `GET /api/game/towns/:id` - 特定の町の詳細取得
- `GET /api/game/ores` - 鉱石の一覧取得
- `GET /api/game/ores/:id` - 特定の鉱石の詳細取得

### 管理者関連

- `POST /api/admin/login` - 管理者ログイン
- `GET /api/admin/users` - ユーザー一覧取得
- `GET /api/admin/users/:id` - 特定のユーザー詳細取得
- `PUT /api/admin/users/:id` - ユーザー情報更新
- `DELETE /api/admin/users/:id` - ユーザー削除
- `GET /api/admin/pending-users` - 未確認ユーザー一覧取得
- `DELETE /api/admin/pending-users/:token` - 未確認ユーザー削除

## よくある開発タスク

### 新しいAPIエンドポイントの追加

1. `backend/handlers/` に適切なハンドラー関数を追加
2. `backend/simple_server.go` のルーター設定に新しいエンドポイントを追加
3. 必要に応じてモデルや機能を実装
4. フロントエンドのサービスレイヤーに対応するAPI呼び出しを追加

### 新しいビューの追加

1. `frontend/src/views/` に新しいビューコンポーネントを作成
2. `frontend/src/router/index.js` にルート定義を追加
3. 必要に応じてVuexストアを更新
4. ナビゲーションメニューに新しいビューへのリンクを追加
