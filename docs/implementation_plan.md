# 実装計画書

## 概要

DynamoDB と Clean Architecture を組み合わせたオンラインショップ API の段階的実装計画

## 実装方針

### 基本方針

1. **テスト駆動開発**: Red → Green → Refactor サイクル
2. **段階的実装**: 最小単位から徐々に機能拡張
3. **型安全性重視**: コンパイル時エラー検出を最大化
4. **依存性逆転**: Clean Architecture の原則遵守

### 技術スタック確定

```
言語: Go 1.22
Web Framework: Echo v4
DB: DynamoDB (Local)
コード生成: oapi-codegen
テスト: testify + gomock
コンテナ: Docker + docker-compose
AWS SDK: aws-sdk-go-v2
```

## Phase 1: Foundation（目標期間: 1 週間）

### 🎯 目標

プロジェクトの基盤となる環境構築とドメイン層実装

### Sprint 1.1: プロジェクト初期化（2 日目標）

#### Day 1: 環境セットアップ

- [ ] **Task 1.1.1**: Go module 初期化
  ```bash
  go mod init dynamo-modeling
  ```
- [ ] **Task 1.1.2**: プロジェクト構造作成
  ```
  cmd/server/main.go
  internal/domain/entity/
  internal/domain/value/
  internal/domain/repository/
  internal/usecase/
  internal/adapter/controller/
  internal/adapter/presenter/
  internal/adapter/repository/
  internal/handler/
  internal/infrastructure/
  api/openapi.yml
  docker-compose.yml
  Makefile
  ```
- [ ] **Task 1.1.3**: 基本依存関係追加
  ```bash
  go get github.com/labstack/echo/v4
  go get github.com/aws/aws-sdk-go-v2/service/dynamodb
  go get github.com/stretchr/testify
  ```
- [ ] **Task 1.1.4**: Docker Compose 設定
  ```yaml
  services:
    dynamodb-local:
      image: amazon/dynamodb-local:2.4.0
      ports:
        - "8000:8000"
      command: ["-jar", "DynamoDBLocal.jar", "-inMemory", "-port", "8000"]
  ```

#### Day 2: DynamoDB 接続確認

- [ ] **Task 1.1.5**: DynamoDB Local 起動確認
- [ ] **Task 1.1.6**: Go から DynamoDB Local 接続テスト
- [ ] **Task 1.1.7**: テーブル作成スクリプト作成

### Sprint 1.2: ドメイン層実装（3 日目標）

#### Day 3: 値オブジェクト実装

- [ ] **Task 1.2.1**: 基本型定義
  ```go
  type CustomerID string
  type ProductID string
  type OrderID string
  type Money int // cents
  ```
- [ ] **Task 1.2.2**: Email 値オブジェクト
  ```go
  type Email struct { value string }
  func NewEmail(email string) (Email, error)
  ```
- [ ] **Task 1.2.3**: Money 値オブジェクト
  ```go
  type Money struct { cents int }
  func NewMoney(amount int) (Money, error)
  func (m Money) Add(other Money) Money
  ```
- [ ] **Task 1.2.4**: 値オブジェクトの単体テスト

#### Day 4: エンティティ実装

- [ ] **Task 1.2.5**: Customer エンティティ
  ```go
  type Customer struct {
    id CustomerID
    email Email
    name string
    // ...
  }
  func NewCustomer(...) *Customer
  ```
- [ ] **Task 1.2.6**: Product エンティティ
- [ ] **Task 1.2.7**: Order エンティティ（基本構造）
- [ ] **Task 1.2.8**: エンティティの単体テスト

#### Day 5: リポジトリインターフェース

- [ ] **Task 1.2.9**: CustomerRepository インターフェース
  ```go
  type CustomerRepository interface {
    Save(ctx context.Context, customer *Customer) error
    FindByID(ctx context.Context, id CustomerID) (*Customer, error)
    FindByEmail(ctx context.Context, email Email) (*Customer, error)
  }
  ```
- [ ] **Task 1.2.10**: ProductRepository インターフェース
- [ ] **Task 1.2.11**: OrderRepository インターフェース
- [ ] **Task 1.2.12**: インメモリリポジトリ実装（テスト用）

## Phase 2: Infrastructure（目標期間: 1 週間）

### 🎯 目標

DynamoDB 接続とデータ永続化機能の実装

### Sprint 2.1: DynamoDB セットアップ（2 日目標）

#### Day 6: テーブル設計実装

- [ ] **Task 2.1.1**: DynamoDB テーブル作成スクリプト
  ```bash
  aws dynamodb create-table --table-name OnlineShop \
    --attribute-definitions AttributeName=PK,AttributeType=S AttributeName=SK,AttributeType=S \
    --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
    --endpoint-url http://localhost:8000
  ```
- [ ] **Task 2.1.2**: GSI1, GSI2 作成
- [ ] **Task 2.1.3**: テーブル作成の自動化（Makefile）

#### Day 7: AWS SDK 設定

- [ ] **Task 2.1.4**: DynamoDB クライアント設定
  ```go
  cfg, err := config.LoadDefaultConfig(ctx,
    config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
      return aws.Endpoint{URL: "http://localhost:8000"}, nil
    })))
  ```
- [ ] **Task 2.1.5**: 設定の環境変数化
- [ ] **Task 2.1.6**: 接続ヘルスチェック実装

### Sprint 2.2: リポジトリ実装（3 日目標）

#### Day 8: データマッパー実装

- [ ] **Task 2.2.1**: Customer データマッパー
  ```go
  func CustomerToItem(customer *entity.Customer) (map[string]types.AttributeValue, error)
  func ItemToCustomer(item map[string]types.AttributeValue) (*entity.Customer, error)
  ```
- [ ] **Task 2.2.2**: Product データマッパー
- [ ] **Task 2.2.3**: マッパーの単体テスト

#### Day 9: Customer リポジトリ

- [ ] **Task 2.2.4**: DynamoCustomerRepository 実装
  ```go
  func (r *DynamoCustomerRepository) Save(ctx context.Context, customer *entity.Customer) error
  func (r *DynamoCustomerRepository) FindByID(ctx context.Context, id value.CustomerID) (*entity.Customer, error)
  ```
- [ ] **Task 2.2.5**: Customer リポジトリの統合テスト

#### Day 10: Product/Order リポジトリ

- [ ] **Task 2.2.6**: DynamoProductRepository 実装
- [ ] **Task 2.2.7**: DynamoOrderRepository 実装（基本機能）
- [ ] **Task 2.2.8**: 全リポジトリの統合テスト

## Phase 3: API Layer（目標期間: 1 週間）

### 🎯 目標

REST API エンドポイントと業務ユースケースの実装

### Sprint 3.1: OpenAPI & 基本 API（2 日目標）

#### Day 11: OpenAPI 仕様定義

- [ ] **Task 3.1.1**: OpenAPI 仕様書作成
  ```yaml
  openapi: 3.1.0
  info:
    title: Online Shop API
    version: 1.0.0
  paths:
    /customers:
      post: ...
      get: ...
    /products:
      post: ...
      get: ...
  ```
- [ ] **Task 3.1.2**: oapi-codegen 設定
- [ ] **Task 3.1.3**: API 型とサーバーインターフェース生成

#### Day 12: 基本 CRUD API

- [ ] **Task 3.1.4**: Customer CRUD API 実装
  ```go
  func (h *CustomerHandler) PostCustomers(ctx echo.Context) error
  func (h *CustomerHandler) GetCustomer(ctx echo.Context, customerId string) error
  ```
- [ ] **Task 3.1.5**: Product CRUD API 実装
- [ ] **Task 3.1.6**: 基本 API の動作確認

### Sprint 3.2: 業務 API（3 日目標）

#### Day 13: ユースケース層実装

- [ ] **Task 3.2.1**: CreateCustomerUseCase 実装
  ```go
  type CreateCustomerUseCase struct {
    customerRepo repository.CustomerRepository
  }
  func (uc *CreateCustomerUseCase) Execute(ctx context.Context, cmd CreateCustomerCommand) (*Customer, error)
  ```
- [ ] **Task 3.2.2**: CreateProductUseCase 実装
- [ ] **Task 3.2.3**: ユースケースの単体テスト

#### Day 14: 注文機能実装

- [ ] **Task 3.2.4**: CreateOrderUseCase 実装
- [ ] **Task 3.2.5**: GetCustomerOrdersUseCase 実装（GSI2 使用）
- [ ] **Task 3.2.6**: 注文 API エンドポイント実装

#### Day 15: エラーハンドリング・E2E テスト

- [ ] **Task 3.2.7**: 統一エラーハンドリング実装
- [ ] **Task 3.2.8**: バリデーション実装
- [ ] **Task 3.2.9**: E2E テスト実装

## Phase 4: Quality & Documentation（目標期間: 3 日）

### 🎯 目標

品質向上とドキュメント整備

#### Day 16: 品質向上

- [ ] **Task 4.1**: テストカバレッジ向上（目標 80%以上）
- [ ] **Task 4.2**: ログ実装（slog 使用）
- [ ] **Task 4.3**: メトリクス実装（基本的なもの）

#### Day 17: パフォーマンス・セキュリティ

- [ ] **Task 4.4**: DynamoDB クエリ最適化確認
- [ ] **Task 4.5**: セキュリティヘッダー追加
- [ ] **Task 4.6**: レート制限実装（基本的なもの）

#### Day 18: ドキュメント整備

- [ ] **Task 4.7**: API 仕様書完成
- [ ] **Task 4.8**: README 更新（開発環境構築手順）
- [ ] **Task 4.9**: 学習振り返りドキュメント作成

## 進捗管理

### デイリーチェックポイント

- [ ] その日のタスクが完了したか
- [ ] テストが通っているか
- [ ] 設計原則から逸脱していないか
- [ ] 翌日のタスクが明確か

### ウィークリーレビュー

- [ ] Phase 目標達成度
- [ ] 学習内容の振り返り
- [ ] 次 Phase の準備確認
- [ ] リスクアセスメント更新

### 完了条件

1. ✅ 全機能が DynamoDB Local で動作
2. ✅ テストカバレッジ 80%以上
3. ✅ OpenAPI 仕様書完成
4. ✅ Clean Architecture の原則遵守
5. ✅ 学習目標達成（アクセスパターン理解等）

## 緊急時対応

### スケジュール遅延時

1. **1 日遅延**: タスクの優先度見直し、不要機能の削除
2. **3 日遅延**: Phase scope の縮小検討
3. **1 週間遅延**: MVP 再定義、Phase 4 の簡素化

### 技術的課題発生時

1. **DynamoDB 設計問題**: 一時的に単純な設計に変更
2. **Go 言語理解不足**: より簡単な実装方法に変更
3. **Clean Architecture 理解不足**: レイヤー構造の簡素化

### リソース不足時

- **時間不足**: 自動テストを手動テストに変更
- **知識不足**: 公式ドキュメント・サンプルコード優先参照
- **環境不足**: クラウド環境の利用検討
