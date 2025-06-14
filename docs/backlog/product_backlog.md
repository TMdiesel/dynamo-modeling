# プロダクトバックログ

## MVP（Minimum Viable Product）

### 目標

DynamoDB の Single Table Design を学習するための最小限のオンラインショップ API

### MVP スコープ

- 顧客管理（登録・取得）
- 商品管理（登録・取得・一覧）
- 注文管理（作成・取得・顧客別履歴）
- DynamoDB Local 環境での動作確認

## ユーザーストーリー

### Epic 1: 環境構築とプロジェクト初期化

- **US-001**: 開発者として、Go module とプロジェクト構造を初期化したい
- **US-002**: 開発者として、DynamoDB Local を docker-compose で起動したい
- **US-003**: 開発者として、OpenAPI 仕様を定義してコード生成したい

### Epic 2: ドメイン層実装

- **US-101**: 開発者として、値オブジェクト（CustomerID, Email, Money 等）を実装したい
- **US-102**: 開発者として、ドメインエンティティ（Customer, Product, Order）を実装したい
- **US-103**: 開発者として、リポジトリインターフェースを定義したい

### Epic 3: インフラ層実装（DynamoDB）

- **US-201**: 開発者として、DynamoDB テーブルを作成したい
- **US-202**: 開発者として、Customer 用のリポジトリを実装したい
- **US-203**: 開発者として、Product 用のリポジトリを実装したい
- **US-204**: 開発者として、Order 用のリポジトリを実装したい

### Epic 4: API 層実装

- **US-301**: 開発者として、Customer API エンドポイントを実装したい
- **US-302**: 開発者として、Product API エンドポイントを実装したい
- **US-303**: 開発者として、Order API エンドポイントを実装したい

### Epic 5: テストと品質保証

- **US-401**: 開発者として、ドメインロジックの単体テストを書きたい
- **US-402**: 開発者として、リポジトリの統合テストを書きたい
- **US-403**: 開発者として、API の E2E テストを書きたい

## リリース計画

### Phase 1: Foundation（1 週間）

**目標**: プロジェクトの土台となる環境とドメイン層を構築

#### Sprint 1.1: プロジェクト初期化（2 日）

- Go module 初期化
- プロジェクト構造作成
- Docker Compose 環境構築
- DynamoDB Local 接続確認

#### Sprint 1.2: ドメイン層実装（3 日）

- 値オブジェクトの実装
- エンティティの実装
- リポジトリインターフェース定義
- ドメインロジックの単体テスト

### Phase 2: Infrastructure（1 週間）

**目標**: DynamoDB 接続とデータ永続化機能を実装

#### Sprint 2.1: DynamoDB セットアップ（2 日）

- テーブル作成スクリプト
- AWS SDK 設定
- 接続テスト

#### Sprint 2.2: リポジトリ実装（3 日）

- Customer リポジトリ実装
- Product リポジトリ実装
- Order リポジトリ実装
- 統合テスト

### Phase 3: API Layer（1 週間）

**目標**: REST API エンドポイントと業務機能を実装

#### Sprint 3.1: OpenAPI & 基本 API（2 日）

- OpenAPI 仕様定義
- コード生成セットアップ
- 基本的な CRUD API 実装

#### Sprint 3.2: 業務 API（3 日）

- 注文作成フロー
- 顧客別注文履歴 API
- エラーハンドリング
- E2E テスト

### Phase 4: Quality & Documentation（3 日）

**目標**: 品質向上とドキュメント整備

- パフォーマンステスト
- セキュリティ考慮
- API 仕様書完成
- 学習ドキュメント更新

## 優先度マトリクス

| Epic       | ビジネス価値 | 技術的重要度 | 学習効果 | 実装難易度 | 優先度  |
| ---------- | ------------ | ------------ | -------- | ---------- | ------- |
| 環境構築   | 低           | 高           | 中       | 低         | 🔴 最高 |
| ドメイン層 | 中           | 高           | 高       | 中         | 🔴 最高 |
| インフラ層 | 中           | 高           | 高       | 高         | 🟡 高   |
| API 層     | 高           | 中           | 中       | 中         | 🟡 高   |
| テスト     | 中           | 高           | 高       | 中         | 🟢 中   |

## 成功指標

### 技術指標

- [ ] DynamoDB Local 環境でアプリケーションが動作する
- [ ] すべての API エンドポイントが正常に応答する
- [ ] 単体テスト カバレッジ > 80%
- [ ] 統合テスト がすべて成功する

### 学習指標

- [ ] Single Table Design のアクセスパターンを理解している
- [ ] Clean Architecture の各層の責務を説明できる
- [ ] DynamoDB の GSI 使用方法を理解している
- [ ] Go での型安全なドメインモデリングができる

## リスク管理

### 高リスク

- **DynamoDB 設計の複雑性**: アクセスパターンが複雑になりすぎる可能性
  - 対策: MVP では基本的なパターンのみに限定

### 中リスク

- **Go 言語の習熟度**: generics や高度な機能の理解不足
  - 対策: シンプルな実装から始めて段階的に高度化

### 低リスク

- **AWS SDK 使用方法**: ドキュメント豊富で解決しやすい
  - 対策: 公式ドキュメントとサンプルコード参照
