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
  username VARCHAR(255) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  is_admin BOOLEAN DEFAULT FALSE,
  is_email_verified BOOLEAN DEFAULT FALSE,
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

#### player_bases テーブル

プレイヤーの拠点情報を管理します。

```sql
CREATE TABLE player_bases (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  town_id VARCHAR(36) NOT NULL,
  level INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (town_id) REFERENCES towns(id) ON DELETE CASCADE
);
```

#### player_inventories テーブル

プレイヤーの所持金と在庫容量を管理します。

```sql
CREATE TABLE player_inventories (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  gold INT NOT NULL DEFAULT 0,
  max_capacity INT NOT NULL DEFAULT 100,
  current_capacity INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

#### player_ores テーブル

プレイヤーが所持している鉱石の数量を管理します。

```sql
CREATE TABLE player_ores (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  ore_id VARCHAR(36) NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (ore_id) REFERENCES ores(id) ON DELETE CASCADE
);
```

#### player_items テーブル

プレイヤーが所持しているアイテムの数量を管理します。

```sql
CREATE TABLE player_items (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  item_id VARCHAR(36) NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
```

### リレーションシップ

- ユーザー (1) -> (N) 拠点 (`player_bases`)
- 町 (1) -> (N) 拠点 (`player_bases`)
- ユーザー (1) -> (1) プレイヤーインベントリ (`player_inventories`)
- ユーザー (1) -> (N) プレイヤー所持鉱石 (`player_ores`)
- 鉱石 (1) -> (N) プレイヤー所持鉱石 (`player_ores`)
- ユーザー (1) -> (N) プレイヤー所持アイテム (`player_items`)
- アイテム (1) -> (N) プレイヤー所持アイテム (`player_items`)

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

データベーススキーマの変更は、バックエンドの起動時に自動的に適用されます。実際の環境では適切なマイグレーションツールの導入を検討する必要があります。

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
