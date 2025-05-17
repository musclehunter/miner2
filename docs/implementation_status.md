# 「採掘人」ゲーム実装状況

## プロジェクト概要

「採掘人」は、プレイヤーが鉱石を採掘して売却するシミュレーションゲームです。

### 技術スタック

- **サーバーサイド**: Go言語 + Gin Webフレームワーク
- **クライアントサイド**: Vue.js
- **データストア**: MySQL, Redis
- **開発環境**: Docker Compose

## 現在の実装状況

### 1. 認証システム

認証システムは完全に実装され、以下の機能を提供しています：

- **ユーザー登録 (Signup)**: 新規ユーザーの作成
- **ログイン (Login)**: 既存ユーザーの認証
- **認証ミドルウェア**: 保護されたエンドポイントへのアクセス制御

ユーザー認証には以下のセキュリティ対策が実装されています：
- パスワードのbcryptハッシュ化（ソルト付き）
- JWTトークンによるセッション管理

### 2. ゲームデータ基盤

ゲームの基本データは以下のモデルで管理されています：

#### 町 (Town)

町はプレイヤーが訪れて鉱石を採掘できる場所です。

```go
type Town struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### 鉱石 (Ore)

鉱石はゲーム内で採掘・取引できる資源です。

```go
type Ore struct {
    ID                  string    `json:"id" gorm:"primaryKey"`
    Name                string    `json:"name" gorm:"not null"`
    Rarity              int       `json:"rarity" gorm:"not null"`
    Purity              float64   `json:"purity" gorm:"not null"`
    ProcessingDifficulty int       `json:"processing_difficulty" gorm:"not null"`
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
}
```

### 3. APIエンドポイント

現在、以下のAPIエンドポイントが実装されています：

#### 認証関連

- `POST /api/auth/signup`: 新規ユーザー登録
- `POST /api/auth/login`: ユーザーログイン
- `GET /api/auth/me`: 現在ログイン中のユーザー情報取得（認証必須）

#### ゲーム関連

- `GET /api/game/towns`: すべての町情報を取得
- `GET /api/game/towns/:id`: 指定IDの町情報を取得
- `GET /api/game/ores`: すべての鉱石情報を取得
- `GET /api/game/ores/:id`: 指定IDの鉱石情報を取得

### 4. 初期データ

ゲーム起動時に以下の初期データが自動的に作成されます：

#### 町データ

- アイアンヒル：鉄鉱石の産地として知られる古い鉱山の町
- シルバーレイク：銀鉱石が豊富な湖のほとりにある町
- ゴールドクレスト：金鉱脈が発見されて栄えた歴史ある町
- クリスタルヴェイル：美しい結晶が取れる渓谷近くの町
- コッパークリーク：銅鉱石の採掘で栄えた小さな町

#### 鉱石データ

| ID | 名前 | レア度 | 純度 | 加工難易度 |
|----|------|--------|------|------------|
| 1 | 石炭 | 1 | 1.0 | 1 |
| 2 | 鉄鉱石 | 1 | 1.0 | 2 |
| 3 | 銅鉱石 | 2 | 1.0 | 3 |
| 4 | 銀鉱石 | 3 | 0.9 | 4 |
| 5 | 金鉱石 | 4 | 0.8 | 5 |
| 6 | ダイヤモンド原石 | 5 | 0.7 | 7 |
| 7 | エメラルド原石 | 5 | 0.7 | 6 |
| 8 | サファイア原石 | 5 | 0.7 | 6 |
| 9 | ルビー原石 | 5 | 0.7 | 6 |
| 10 | ミスリル鉱石 | 6 | 0.5 | 10 |

## データベース構造

主要なテーブル構造は以下の通りです：

### usersテーブル

```sql
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### townsテーブル

```sql
CREATE TABLE IF NOT EXISTS towns (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### oresテーブル

```sql
CREATE TABLE IF NOT EXISTS ores (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rarity INT NOT NULL,
    purity FLOAT NOT NULL,
    processing_difficulty INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 今後の開発予定

### 1. プレイヤー固有のゲームデータ管理

- プレイヤーの所持金や鉱石在庫などの管理機能
- 採掘メカニズムの実装

### 2. ゲームロジックの拡張

- 町ごとに異なる鉱石の出現率や特性の設定
- 採掘の成功率や採掘量の計算ロジック

### 3. フロントエンド連携

- 認証機能と町・鉱石情報の表示機能をフロントエンドと連携
- ゲームUIの実装

## 実行方法

### 開発環境の起動

```bash
# Docker環境を起動
docker-compose up -d

# バックエンドのみ再起動
docker-compose restart backend
```

### APIアクセス例

#### 町情報の取得

```bash
curl http://localhost:8080/api/game/towns
```

#### 鉱石情報の取得

```bash
curl http://localhost:8080/api/game/ores
```
