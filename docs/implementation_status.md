# 採掘人 - 実装状況

## 概要

このドキュメントでは「採掘人」ゲームの実装状況を記録します。

## フロントエンド

### 実装済み画面
- ✅ ログイン画面
- ✅ ユーザー登録画面
- ✅ メインダッシュボード
- ✅ ワールドマップ（座標ベースの町選択機能付き）
- ✅ 町詳細画面

### 進行中の実装
- 🔄 インベントリ管理画面
- 🔄 鉱山掘削画面
- 🔄 労働者管理画面

### 未実装
- ❌ マーケット取引画面
- ❌ クエスト管理画面
- ❌ ユーザープロフィール画面

## バックエンド

### 実装済みAPI
- ✅ 認証API（ログイン・登録）
- ✅ 町情報取得API
- ✅ 鉱石情報API
- ✅ プレイヤー在庫情報取得API (`/api/game/my/inventory`)

### 進行中の実装
- 🔄 拠点管理API

### 未実装
- ❌ 取引API
- ❌ クエストAPI
- ❌ 労働者管理API

## データベース

### 実装済みテーブル
- ✅ ユーザーテーブル
- ✅ 町テーブル（座標情報を含む）
- ✅ 鉱石テーブル (定義を更新)
- ✅ アイテムテーブル (新規追加)
- ✅ プレイヤーインベントリテーブル (`player_inventories`)
- ✅ プレイヤー所持鉱石テーブル (`player_ores`)
- ✅ プレイヤー所持アイテムテーブル (`player_items`)

### 進行中のテーブル実装
- 🔄 労働者テーブル

## 最近の更新

### 2025-06-11
- **プレイヤー在庫システムのバックエンド実装完了**
  - モデル、リポジトリ、ハンドラーを`database/sql`ベースで実装完了
  - GORMへの依存を排除
  - 認証済みユーザーの在庫情報をまとめて取得するAPIエンドポイント `GET /api/game/my/inventory` を実装
- プレイヤーのインベントリシステム実装開始
  - データベース設計更新:
    - `items`テーブルを新規設計
    - `ores`テーブル定義を`project_config.json`に合わせて精緻化
    - `player_inventories`テーブル（所持金・在庫容量）を設計
    - プレイヤー所持リソースの管理を`player_ores`（所持鉱石）と`player_items`（所持アイテム）テーブルに分割して設計
  - APIエンドポイントの設計変更:
    - プレイヤー所持金・在庫取得API (`/api/game/inventory`) を定義
    - プレイヤー所持鉱石取得API (`/api/game/player-ores`) を定義
    - プレイヤー所持アイテム取得API (`/api/game/player-items`) を定義
    - (旧 `/api/game/resources` は上記2つに分割)
  - プロジェクト設定ファイル（`.windsurf/project_config.json`）のデータモデルとAPIエンドポイント定義を上記変更に合わせて更新
  - データベース設計ドキュメント (`database.md`) を上記変更に合わせて更新

## プレイヤー在庫システム - バックエンド実装タスク

### 1. モデル定義 (Go structs in `backend/models/`)
- [x] `player_inventory.go` (for `PlayerInventory` model)
- [x] `player_ore.go` (for `PlayerOre` model)
- [x] `player_item.go` (for `PlayerItem` model)

### 2. データベースリポジトリ実装 (CRUD operations in `backend/database/`)
- [x] `player_inventory_repository.go`
- [x] `player_ore_repository.go`
- [x] `player_item_repository.go`

### 3. APIハンドラー実装 (in `backend/handlers/`)
- [x] `inventory_handler.go`を作成し、`GetMyInventory`ハンドラーを実装

### 4. APIルーティング (in Gin router setup)
- [x] `/api/game/my/inventory` のルートを登録し、認証ミドルウェアを適用

## プレイヤー在庫システム - フロントエンド実装タスク

### 1. APIサービス更新 (in `frontend/src/services/`)
- [ ] `api.js` (または新規ファイル) に、`/api/game/inventory`, `/api/game/player-ores`, `/api/game/player-items` を呼び出す関数を追加。

### 2. Vuexストア更新 (in `frontend/src/store/`)
- [ ] プレイヤーの所持金、鉱石リスト、アイテムリストを管理するためのstate, getters, mutations, actionsを追加。

### 3. UIコンポーネント作成 (in `frontend/src/components/` or `frontend/src/views/`)
- [ ] プレイヤー所持金表示用コンポーネント (`PlayerGoldDisplay.vue` など)
- [ ] プレイヤー所持鉱石リスト表示用コンポーネント (`PlayerOreList.vue` など)
- [ ] プレイヤー所持アイテムリスト表示用コンポーネント (`PlayerItemList.vue` など)

### 4. 拠点画面への統合 (e.g., `BaseScreen.vue`)
- [ ] 上記で作成したコンポーネントを拠点画面に配置し、APIから取得したデータを表示。

### 2025-06-02
- ワールドマップ画面の刷新
  - 新しいマップ画像を使用
  - 町の座標ベースの選択機能を実装
  - マーカー表示とホバーエフェクト追加
- 町データの更新
  - 7つの地域に対応する町データを実装
  - `world_setting.md`に基づく詳細な世界観を反映
- データベースマイグレーション実施
  - 町テーブルに座標カラム（position_x, position_y）を追加
  - 既存の町データに座標情報を設定
  - 不足していた2つの町（フォグヴェイル、キャメロス）を追加
- データベーススキーマの更新
  - townテーブルに座標情報（position_x, position_y）を追加
- バックエンドリポジトリの更新
  - Town関連の全メソッド（GetAllTowns, GetTownByID, CreateTown, UpdateTown）を座標対応に更新
- フロントエンドの改善
  - 7つの地域に合わせたオリジナルSVGアイコンを作成・実装
  - ワールドマップでの町選択エラーを修正
- APIドキュメント更新
  - 町データのフィールド名を実際の実装（position_x, position_y）に合わせて更新
