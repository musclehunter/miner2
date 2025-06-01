# 採掘人（Miner2）

鉱石採掘シミュレーションゲーム「採掘人」のプロジェクトリポジトリへようこそ。

## プロジェクト概要

「採掘人」は、プレイヤーが鉱石を採掘して売却するシミュレーションゲームです。Go言語のバックエンドとVue.jsのフロントエンドで構成されています。

### 技術スタック

- **サーバーサイド**: Go言語 + Gin Webフレームワーク
- **クライアントサイド**: Vue.js
- **データストア**: MySQL, Redis
- **開発環境**: Docker Compose

## クイックスタート

```bash
# リポジトリのクローン
git clone <repository-url>
cd miner2

# Dockerコンテナの起動
docker-compose up -d

# フロントエンドへのアクセス
open http://localhost:8081
```

## 開発ドキュメント

詳細なドキュメントは以下のリンクから参照できます：

- [インストールガイド](docs/installation.md) - セットアップ手順と環境構成
- [システムアーキテクチャ](docs/architecture.md) - システム全体の設計と構成
- [開発ガイド](docs/development.md) - 開発環境と開発フロー
- [データベース設計](docs/database.md) - データベーススキーマと構造
- [API仕様書](docs/api.md) - RESTful APIエンドポイントの詳細
- [実装状況](docs/implementation_status.md) - 現在の実装状況と今後の計画

## 機能ハイライト

- **ユーザー認証システム** - サインアップ、ログイン、メール確認
- **ゲームデータ管理** - 町、鉱石、プレイヤーデータの管理
- **管理者インターフェース** - ユーザーとゲームデータの管理
- **グローバルナビゲーション** - すべての画面で一貫したUI

## 貢献ガイドライン

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを開く

## ライセンス

このプロジェクトはプライベートリポジトリです。すべての権利は保持されています。

## 連絡先

プロジェクト管理者にお問い合わせください。
