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
- ✅ 拠点画面

### 進行中の実装
- ✅ 拠点設立機能
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
- ✅ 管理者向け拠点情報取得API (`/api/admin/bases`)

### 進行中の実装
- ✅ 拠点設立API

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
- ✅ プレイヤー拠点テーブル (`player_bases`)

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
- [x] `frontend/src/services/inventoryService.js` を作成し、`/api/game/my/inventory` を呼び出す関数を追加。

### 2. Vuexストア更新 (in `frontend/src/store/`)
- [x] プレイヤーの所持金、鉱石リスト、アイテムリストを管理するためのVuexモジュール (`frontend/src/store/modules/inventory.js`) を作成し、state, getters, mutations, actionsを追加。

### 3. UIコンポーネント作成 (in `frontend/src/components/` or `frontend/src/views/`)
- [x] プレイヤー所持金表示用コンポーネント (`frontend/src/components/inventory/PlayerGoldDisplay.vue`) を作成。
- [x] プレイヤー所持鉱石リスト表示用コンポーネント (`frontend/src/components/inventory/PlayerOreList.vue`) を作成。
- [x] プレイヤー所持アイテムリスト表示用コンポーネント (`frontend/src/components/inventory/PlayerItemList.vue`) を作成。

### 4. 拠点画面への統合 (e.g., `BaseScreen.vue`)
- [x] 上記で作成したコンポーネントを拠点画面 (`frontend/src/views/BaseView.vue`) に配置し、APIから取得したデータを表示。

## 拠点設立システム - 実装タスク

### バックエンド実装タスク

1.  **モデル定義 (in `backend/models/`)**
    - [x] `player_base.go` (for `PlayerBase` model) - `database.md`を参考に既存定義を確認し、必要に応じて修正。

2.  **データベースリポジトリ実装 (in `backend/database/`)**
    - [x] `player_base_repository.go` を作成し、拠点を新規作成する `CreatePlayerBase` 関数を実装。この関数はトランザクション内で `player_inventories` テーブルにも初期レコードを作成する。

3.  **APIハンドラー実装 (in `backend/handlers/`)**
    - [x] `base_handler.go` を作成し、`CreateBase` ハンドラーを実装。リクエストから町IDを受け取り、リポジトリを呼び出す。

4.  **APIルーティング (in Gin router setup)**
    - [x] `POST /api/game/bases` のルートを登録し、認証ミドルウェアを適用。

### フロントエンド実装タスク

1.  **APIサービス更新 (in `frontend/src/services/`)**
    - [x] `baseService.js` を新規作成し、`POST /api/game/bases` を呼び出す `createBase` 関数を追加。

2.  **Vuexストア更新 (in `frontend/src/store/`)**
    - [x] プレイヤーの拠点情報を管理するための `base` モジュールを `modules` に追加。拠点の作成、取得に関する state, actions を定義。

3.  **UIコンポーネント/ビュー更新**
    - [x] `WorldMapView.vue` を改修。町を選択した際に、プレイヤーがまだその町に拠点を持っていない場合、拠点設立の確認モーダルと「設立」ボタンを表示する。
    - [x] 拠点設立の確認モーダルで「設立」ボタンクリック時に `createBase` アクションをディスパッチする。
    - [x] 拠点設立成功後、ユーザーに通知し、拠点画面へ遷移させる。

## 管理機能

- ✅ ユーザー管理
- ✅ 未確認ユーザー管理
- ✅ 町管理
- ✅ 拠点管理

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
