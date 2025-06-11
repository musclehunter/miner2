# 採掘人 - API仕様書

## 概要

「採掘人」のバックエンドAPIはREST形式で提供されており、JSONデータの送受信を行います。

## 認証

保護されたエンドポイントにアクセスするには、認証が必要です。認証はJWTトークンを使用します。

### 認証ヘッダー

```
Authorization: Bearer <jwt_token>
```

## エンドポイント

### 認証関連 API

#### ユーザー登録

新規ユーザーを登録し、確認メールを送信します。

- **URL**: `/api/auth/signup`
- **Method**: `POST`
- **認証**: 不要
- **リクエスト**:

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "name": "Test User"
}
```

- **レスポンス**:

```json
{
  "message": "確認メールを送信しました。メールを確認して登録を完了してください。",
  "email": "test@example.com"
}
```

#### ログイン

登録済みユーザーの認証を行い、JWTトークンを返します。

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **認証**: 不要
- **リクエスト**:

```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

- **レスポンス**:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "testuser",
    "email": "test@example.com",
    "name": "Test User",
    "is_admin": false
  }
}
```

#### 現在のユーザー情報取得

認証されたユーザーの情報を取得します。

- **URL**: `/api/auth/me`
- **Method**: `GET`
- **認証**: 必要
- **レスポンス**:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "testuser",
  "email": "test@example.com",
  "name": "Test User",
  "is_admin": false
}
```

#### メールアドレス確認

登録時に送信された確認トークンを検証します。

- **URL**: `/api/auth/verify-email`
- **Method**: `GET`
- **認証**: 不要
- **クエリパラメータ**: `token=<verification_token>`
- **レスポンス**:

```json
{
  "message": "メールアドレスが確認されました。"
}
```

#### 確認メール再送信

確認メールを再送信します。

- **URL**: `/api/auth/resend-verification`
- **Method**: `POST`
- **認証**: 不要
- **リクエスト**:

```json
{
  "email": "test@example.com"
}
```

- **レスポンス**:

```json
{
  "message": "確認メールを再送信しました。",
  "email": "test@example.com"
}
```

### ゲーム関連 API

#### 町の一覧取得

ゲーム内の町の一覧を取得します。

- **URL**: `/api/game/towns`
- **Method**: `GET`
- **認証**: 必要
- **レスポンス**:

```json
[
  {
    "id": "1",
    "name": "スタートタウン",
    "description": "初心者向けの町",
    "population": 1000,
    "prosperity": 70,
    "security": 80,
    "position_x": 50,
    "position_y": 50
  },
  {
    "id": "2",
    "name": "マウンテンビレッジ",
    "description": "山岳地帯の町",
    "population": 500,
    "prosperity": 40,
    "security": 60,
    "position_x": 80,
    "position_y": 30
  }
]
```

#### 特定の町の詳細取得

特定の町の詳細情報を取得します。

- **URL**: `/api/game/towns/:id`
- **Method**: `GET`
- **認証**: 必要
- **レスポンス**:

```json
{
  "id": "1",
  "name": "スタートタウン",
  "description": "初心者向けの町",
  "population": 1000,
  "prosperity": 70,
  "security": 80,
  "position_x": 50,
  "position_y": 50,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

#### プレイヤー在庫情報取得

ログインしているプレイヤーの在庫情報（所持金、鉱石、アイテム）をまとめて取得します。

- **URL**: `/api/game/my/inventory`
- **Method**: `GET`
- **認証**: 必要
- **レスポンス**:

```json
{
  "inventory": {
    "id": "inv-user-123",
    "user_id": "user-123",
    "gold": 1000,
    "max_capacity": 500,
    "current_capacity": 120,
    "created_at": "2023-01-10T10:00:00Z",
    "updated_at": "2023-01-10T10:00:00Z"
  },
  "ores": [
    {
      "id": "player-ore-1",
      "user_id": "user-123",
      "ore_id": "ore-1",
      "quantity": 50,
      "created_at": "2023-01-10T10:00:00Z",
      "updated_at": "2023-01-10T10:00:00Z",
      "ore": {
        "id": "ore-1",
        "name": "鉄鉱石",
        "rarity": 1,
        "purity": 60,
        "processing_difficulty": 20,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    }
  ],
  "items": [
    {
      "id": "player-item-1",
      "user_id": "user-123",
      "item_id": "item-1",
      "quantity": 5,
      "created_at": "2023-01-10T10:00:00Z",
      "updated_at": "2023-01-10T10:00:00Z",
      "item": {
        "id": "item-1",
        "name": "ピッケル",
        "rarity": 2,
        "description": "頑丈なピッケル。",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    }
  ]
}
```

#### 鉱石の一覧取得

ゲーム内の鉱石の一覧を取得します。

- **URL**: `/api/game/ores`
- **Method**: `GET`
- **認証**: 必要
- **レスポンス**:

```json
[
  {
    "id": "1",
    "name": "銅鉱石",
    "description": "一般的な鉱石",
    "base_value": 10,
    "rarity": 1
  },
  {
    "id": "2",
    "name": "鉄鉱石",
    "description": "基本的な金属鉱石",
    "base_value": 20,
    "rarity": 2
  }
]
```

#### 特定の鉱石の詳細取得

特定の鉱石の詳細情報を取得します。

- **URL**: `/api/game/ores/:id`
- **Method**: `GET`
- **認証**: 必要
- **レスポンス**:

```json
{
  "id": "1",
  "name": "銅鉱石",
  "description": "一般的な鉱石",
  "base_value": 10,
  "rarity": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 管理者関連 API

#### 全ての拠点情報取得

システムに登録されている全てのプレイヤー拠点情報を取得します。

- **URL**: `/api/admin/bases`
- **Method**: `GET`
- **認証**: 必要 (管理者)
- **レスポンス**:

```json
[
  {
    "id": "b8e1b4e6-3e4f-4d2a-8c1a-9b8e1b4e63e4",
    "user_id": "u1",
    "town_id": "t1",
    "level": 1,
    "created_at": "2025-06-12T00:36:06Z",
    "updated_at": "2025-06-12T00:36:06Z"
  },
  {
    "id": "c9f2c5f7-4f5g-5e3b-9d2b-1c9f2c5f74f5",
    "user_id": "u2",
    "town_id": "t2",
    "level": 2,
    "created_at": "2025-06-12T01:00:00Z",
    "updated_at": "2025-06-12T01:30:00Z"
  }
]
```

#### 管理者ログイン

管理者としてログインします。

- **URL**: `/api/admin/login`
- **Method**: `POST`
- **認証**: 不要
- **リクエスト**:

```json
{
  "email": "admin@example.com",
  "password": "adminpassword"
}
```

- **レスポンス**:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "admin",
    "email": "admin@example.com",
    "name": "Admin User",
    "is_admin": true
  }
}
```

#### ユーザー一覧取得

すべてのユーザーの一覧を取得します（管理者のみ）。

- **URL**: `/api/admin/users`
- **Method**: `GET`
- **認証**: 必要（管理者権限）
- **レスポンス**:

```json
[
  {
    "id": "1",
    "username": "user1",
    "email": "user1@example.com",
    "name": "User One",
    "is_admin": false,
    "is_email_verified": true,
    "created_at": "2023-01-01T00:00:00Z"
  },
  {
    "id": "2",
    "username": "user2",
    "email": "user2@example.com",
    "name": "User Two",
    "is_admin": false,
    "is_email_verified": false,
    "created_at": "2023-01-02T00:00:00Z"
  }
]
```

#### 特定のユーザー詳細取得

特定のユーザーの詳細情報を取得します（管理者のみ）。

- **URL**: `/api/admin/users/:id`
- **Method**: `GET`
- **認証**: 必要（管理者権限）
- **レスポンス**:

```json
{
  "id": "1",
  "username": "user1",
  "email": "user1@example.com",
  "name": "User One",
  "is_admin": false,
  "is_email_verified": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

#### ユーザー情報更新

特定のユーザーの情報を更新します（管理者のみ）。

- **URL**: `/api/admin/users/:id`
- **Method**: `PUT`
- **認証**: 必要（管理者権限）
- **リクエスト**:

```json
{
  "name": "Updated User Name",
  "is_admin": true
}
```

- **レスポンス**:

```json
{
  "message": "ユーザー情報が更新されました。",
  "user": {
    "id": "1",
    "username": "user1",
    "email": "user1@example.com",
    "name": "Updated User Name",
    "is_admin": true,
    "is_email_verified": true
  }
}
```

#### ユーザー削除

特定のユーザーを削除します（管理者のみ）。

- **URL**: `/api/admin/users/:id`
- **Method**: `DELETE`
- **認証**: 必要（管理者権限）
- **レスポンス**:

```json
{
  "message": "ユーザーが削除されました。"
}
```

#### 未確認ユーザー一覧取得

メール確認が完了していない仮登録ユーザーの一覧を取得します（管理者のみ）。

- **URL**: `/api/admin/pending-users`
- **Method**: `GET`
- **認証**: 必要（管理者権限）
- **レスポンス**:

```json
[
  {
    "email": "pending1@example.com",
    "name": "Pending User One",
    "created_at": "2023-01-01T00:00:00Z",
    "token": "token1"
  },
  {
    "email": "pending2@example.com",
    "name": "Pending User Two",
    "created_at": "2023-01-02T00:00:00Z",
    "token": "token2"
  }
]
```

#### 未確認ユーザー削除

特定の未確認ユーザーを削除します（管理者のみ）。

- **URL**: `/api/admin/pending-users/:token`
- **Method**: `DELETE`
- **認証**: 必要（管理者権限）
- **レスポンス**:

```json
{
  "message": "未確認ユーザーが削除されました。"
}
```

### 町管理

#### 町一覧取得

全ての町の情報を取得します（管理者のみ）。

- **URL**: `/api/admin/towns`
- **Method**: `GET`
- **認証**: 必要（管理者権限）
- **レスポンス**: `Array of Town objects`

#### 新規町作成

新しい町を作成します（管理者のみ）。

- **URL**: `/api/admin/towns`
- **Method**: `POST`
- **認証**: 必要（管理者権限）
- **リクエスト**: 

```json
{
    "name": "新しい町",
    "description": "説明文"
}
```

- **レスポンス**: `Town object`

#### 町情報更新

- **URL**: `/api/admin/towns/:id`
- **Method**: `PUT`
- **認証**: 必要（管理者権限）
- **リクエスト**: 

```json
{
    "name": "更新された町名",
    "description": "更新された説明文"
}
```

- **レスポンス**: `Town object`

#### 町削除

- **URL**: `/api/admin/towns/:id`
- **Method**: `DELETE`
- **認証**: 必要（管理者権限）
- **レスポンス**: 

```json
{
  "message": "町を削除しました。"
}
```

### 拠点管理

#### 拠点一覧取得

全てのプレイヤー拠点情報を取得します（管理者のみ）。

- **URL**: `/api/admin/bases`
- **Method**: `GET`
- **認証**: 必要（管理者権限）
- **レスポンス**: `Array of PlayerBase objects`


## エラーレスポンス

APIエラーは以下の形式で返されます：

```json
{
  "error": "エラーメッセージ",
  "details": "エラーの詳細情報（ある場合）"
}
```

### 一般的なHTTPステータスコード

- `200 OK`: リクエスト成功
- `201 Created`: リソース作成成功
- `400 Bad Request`: リクエスト形式不正
- `401 Unauthorized`: 認証失敗
- `403 Forbidden`: 権限不足
- `404 Not Found`: リソースが見つからない
- `500 Internal Server Error`: サーバー内部エラー
