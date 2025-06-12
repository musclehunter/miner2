# 採掘人 - データベース設計

## 概要

「採掘人」ゲームは以下のデータストアを使用しています：

1. **MySQL** - 永続的なゲームデータの保存
2. **Redis** - 一時データやキャッシュの管理

## MySQL データベース

### テーブル構造

#### users テーブル

ゲームのユーザー情報を管理します。

```sql
CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### towns テーブル

ゲーム内の町（拠点）データを管理します。

```sql
CREATE TABLE towns (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  position_x INT NOT NULL,
  position_y INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**注**: position_xとposition_yはワールドマップ上の座標を表し、各町の位置を特定するために使用されます。

#### ores テーブル

ゲーム内の鉱石データを管理します。

```sql
CREATE TABLE ores (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  rarity INT NOT NULL,
  purity FLOAT,
  processing_difficulty INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### items テーブル

ゲーム内のアイテムデータを管理します。

```sql
CREATE TABLE items (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  rarity INT NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### bases テーブル

プレイヤーの拠点情報を管理します。

```sql
CREATE TABLE bases (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  town_id VARCHAR(36) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (town_id) REFERENCES towns(id) ON DELETE CASCADE
);
```

#### warehouses テーブル

各拠点に紐づく倉庫の情報を管理します。

```sql
CREATE TABLE warehouses (
  id VARCHAR(36) PRIMARY KEY,
  base_id VARCHAR(36) NOT NULL,
  level INT NOT NULL,
  capacity INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (base_id) REFERENCES bases(id) ON DELETE CASCADE
);
```

#### warehouse_items テーブル

倉庫に保管されているアイテムや鉱石の情報を管理します。

```sql
CREATE TABLE warehouse_items (
  id VARCHAR(36) PRIMARY KEY,
  warehouse_id VARCHAR(36) NOT NULL,
  item_id VARCHAR(36) NULL,
  ore_id VARCHAR(36) NULL,
  quantity INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (ore_id) REFERENCES ores(id) ON DELETE CASCADE
);
```

### リレーションシップ

- ユーザー (1) -> (N) 拠点 (`bases`)
- 町 (1) -> (N) 拠点 (`bases`)
- 拠点 (1) -> (1) 倉庫 (`warehouses`)
- 倉庫 (1) -> (N) 倉庫アイテム (`warehouse_items`)
- アイテム (1) -> (N) 倉庫アイテム (`warehouse_items`)
- 鉱石 (1) -> (N) 倉庫アイテム (`warehouse_items`)

## Redis データストア

Redisは主に一時データの保存やキャッシュとして使用されます。

### キー構造

#### メール確認トークン

```
email_verification:<token> -> <ユーザー情報のJSON>
```

例：
```
email_verification:Wf57KhvLq4LnPLX-srrFnQ9O2c-t9JxB1oz1fBTAP7Q= -> {"email":"test@example.com","name":"Test User","password":"hashedpassword","created_at":"2023-01-01T12:00:00Z"}
```

これらのトークンは24時間後に自動的に期限切れになります。

### Redis設定

Redis接続は以下の環境変数で設定します：

- `REDIS_HOST`: Redisサーバーのホスト名（デフォルト: `redis`）
- `REDIS_PORT`: Redisサーバーのポート（デフォルト: `6379`）
- `REDIS_PASSWORD`: Redisのパスワード（デフォルト: 空）
- `USE_MOCK_REDIS`: モックRedisクライアントの使用有無（true/false）

デフォルトでは `USE_MOCK_REDIS=false` が設定されており、実際のRedisサーバーを使用します。

## データマイグレーション

データベーススキーマの変更は `golang-migrate` を使用して管理されます。
マイグレーションファイルは `database/migrations` ディレクトリに配置されています。
バックエンドの起動時に、保留中のマイグレーションが自動的に適用されます。

## デモデータ

開発環境では初期データとして以下が自動的に投入されます：

- デモユーザー（ユーザー名: `admin`, パスワード: `adminpassword`）
- 複数の町データ
- 複数の鉱石データ

## バックアップと復元

データベースのバックアップを作成するには：

```bash
docker-compose exec db mysqldump -u miner -p minerdb > backup.sql
```

バックアップから復元するには：

```bash
cat backup.sql | docker-compose exec -T db mysql -u miner -p minerdb
