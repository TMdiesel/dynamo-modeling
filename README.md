# DynamoDB + Clean Architecture 学習プロジェクト

DynamoDB のデータモデリングと Go の Clean Architecture を組み合わせた実践的な学習プロジェクトです。

## 🎯 プロジェクト概要

オンラインショップ API を題材として、以下を学習します：

- DynamoDB Single Table Design パターン
- アクセスパターン駆動のデータモデリング
- Go での Clean Architecture 実装
- ドメイン駆動設計（DDD）と NoSQL の組み合わせ

## 🛠 技術スタック

| カテゴリ          | 技術            | バージョン |
| ----------------- | --------------- | ---------- |
| 言語              | Go              | 1.21+      |
| DB                | Amazon DynamoDB | -          |
| Web Framework     | Echo            | v4         |
| AWS SDK           | aws-sdk-go-v2   | latest     |
| Testing           | testify         | v1.8+      |
| API Documentation | OpenAPI 3.0     | -          |

## 🚀 Quick Start

### 前提条件

- Go 1.21 以上
- AWS CLI 設定済み
- Docker（DynamoDB Local 用）

### 開発環境構築

```bash
# リポジトリクローン
git clone <repository-url>
cd dynamo-modeling

# 依存関係インストール
go mod download

# DynamoDB Local起動
docker run -p 8000:8000 amazon/dynamodb-local

# テーブル作成
make setup-local-db

# API起動
make run
```

### API 確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# API仕様確認
open http://localhost:8080/swagger/index.html
```

## 📁 ディレクトリ構成

```
.
├── api/                    # OpenAPI仕様
│   └── openapi.yml
├── cmd/
│   └── server/
│       └── main.go        # アプリケーションエントリーポイント
├── internal/
│   ├── domain/            # ドメイン層（ビジネスロジック）
│   │   ├── entity/        # エンティティ
│   │   ├── value/         # 値オブジェクト
│   │   └── repository/    # リポジトリインターフェース
│   ├── usecase/           # ユースケース層
│   ├── adapter/
│   │   ├── controller/    # コントローラー
│   │   ├── presenter/     # プレゼンター
│   │   ├── repository/    # リポジトリ実装
│   │   └── openapi/       # 生成されたOpenAPI型
│   ├── handler/           # HTTPハンドラー
│   └── infrastructure/    # インフラ層
├── docs/                  # プロジェクトドキュメント
│   ├── strategy/         # ビジョン・ミッション
│   ├── requirements/     # 要件定義
│   ├── design/          # 設計書
│   └── backlog/         # プロダクトバックログ
├── scripts/              # セットアップスクリプト
├── tests/               # テストファイル
└── README.md
```

## 🎯 学習ステップ

### Phase 1: 基礎理解

- [ ] DynamoDB の基本概念
- [ ] Clean Architecture の層構造
- [ ] ドメインエンティティ設計

### Phase 2: 実装基礎

- [ ] Single Table Design 実装
- [ ] 基本 CRUD 操作
- [ ] テストファーストアプローチ

### Phase 3: 応用パターン

- [ ] GSI を使った複雑なクエリ
- [ ] バッチオペレーション
- [ ] パフォーマンス最適化

## 🧪 テスト実行

```bash
# 全テスト実行
make test

# カバレッジ確認
make test-coverage

# 統合テスト（DynamoDB Local必要）
make test-integration
```

## 📚 参考資料

- [AWS DynamoDB Design Patterns](https://github.com/aws-samples/amazon-dynamodb-design-patterns)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [NoSQL Workbench for DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/workbench.html)

## 🔄 ブランチ戦略

- `main`: プロダクション準備済みコード
- `develop`: 開発統合ブランチ
- `feature/*`: 機能開発ブランチ

## 📝 コミット規約

```
feat: 新機能追加
fix: バグ修正
refactor: リファクタリング
test: テスト追加・修正
docs: ドキュメント更新
```

## 🚀 デプロイ

```bash
# AWS環境へのデプロイ
make deploy-staging
make deploy-production
```
