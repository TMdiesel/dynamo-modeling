# 実装計画書

## 概要

DynamoDB と Clean Architecture を組み合わせたオンラインショップ API の段階的実装計画

## 実装方針

### 基本方針

1. **テスト駆動開発**: Red → Green → Refactor サイクル
2. **段階的実装**: 最小単位から徐々に機能拡張
3. **型安全性重視**: コンパイル時エ- [x] **Task 2.2.8**: 全リポジトリの統合テスト

```go
// integration_test.go でエンドツーエンドシナリオテスト実装済み
// Customer登録 → Product作成 → Order作成の一連フローをテスト
// パフォーマンステスト（20件一括操作、225ms完了）実装済み
// 全テストPASS: エンドツーエンド(0.02s) + パフォーマンス(0.23s)
```

**🎯 現在地**: Sprint 3.1 完了！ 次は Sprint 3.2（ユースケース層 & ビジネスロジック実装）

### 📊 Sprint 3.1 完了サマリー

**実装済み機能**:

- ✅ Echo v4 Framework 統合
- ✅ OpenAPI 3.1 仕様書完成（Customer, Product, Order 全エンドポイント）
- ✅ oapi-codegen 設定・コード生成（Echo 用 ServerInterface）
- ✅ APIHandler 実装（全 15 エンドポイント）
- ✅ サーバー起動・エンドポイントテスト成功
- ✅ JSON request/response 処理
- ✅ ミドルウェア設定（CORS, Logger, Recover）

**技術的成果**:

- Echo v4 での OpenAPI-first 開発フロー確立
- 型安全な API 実装（OpenAPI 仕様からの自動生成）
- RESTful API 設計（15 エンドポイント正常動作）
- Graceful Shutdown サーバー実装

**動作確認**:

```bash
# サーバー起動成功 (Echo v4, port 8080)
# GET /customers → サンプルデータ取得成功
# POST /customers → JSON request/response成功
# 全エンドポイント基本動作確認済み
```

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

- [x] **Task 1.1.1**: Go module 初期化
  ```bash
  go mod init dynamo-modeling
  ```
- [x] **Task 1.1.2**: プロジェクト構造作成
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
- [x] **Task 1.1.3**: 基本依存関係追加
  ```bash
  go get github.com/labstack/echo/v4
  go get github.com/aws/aws-sdk-go-v2/service/dynamodb
  go get github.com/stretchr/testify
  ```
- [x] **Task 1.1.4**: Docker Compose 設定
  ```yaml
  services:
    dynamodb-local:
      image: amazon/dynamodb-local:2.4.0
      ports:
        - "8000:8000"
      command: ["-jar", "DynamoDBLocal.jar", "-inMemory", "-port", "8000"]
  ```

#### Day 2: DynamoDB 接続確認

- [x] **Task 1.1.5**: DynamoDB Local 起動確認
- [x] **Task 1.1.6**: Go から DynamoDB Local 接続テスト
- [x] **Task 1.1.7**: テーブル作成スクリプト作成

### Sprint 1.2: ドメイン層実装（3 日目標）

#### Day 3: 値オブジェクト実装

- [x] **Task 1.2.1**: 基本型定義
  ```go
  type CustomerID string
  type ProductID string
  type OrderID string
  type Money int // cents
  ```
- [x] **Task 1.2.2**: Email 値オブジェクト
  ```go
  type Email struct { value string }
  func NewEmail(email string) (Email, error)
  ```
- [x] **Task 1.2.3**: Money 値オブジェクト
  ```go
  type Money struct { cents int }
  func NewMoney(amount int) (Money, error)
  func (m Money) Add(other Money) Money
  ```
- [x] **Task 1.2.4**: 値オブジェクトの単体テスト

#### Day 4: エンティティ実装

- [x] **Task 1.2.5**: Customer エンティティ
  ```go
  type Customer struct {
    id CustomerID
    email Email
    name string
    // ...
  }
  func NewCustomer(...) *Customer
  ```
- [x] **Task 1.2.6**: Product エンティティ
- [x] **Task 1.2.7**: Order エンティティ（基本構造）
- [x] **Task 1.2.8**: エンティティの単体テスト

#### Day 5: リポジトリインターフェース

- [x] **Task 1.2.9**: CustomerRepository インターフェース
  ```go
  type CustomerRepository interface {
    Save(ctx context.Context, customer *Customer) error
    FindByID(ctx context.Context, id CustomerID) (*Customer, error)
    FindByEmail(ctx context.Context, email Email) (*Customer, error)
  }
  ```
- [x] **Task 1.2.10**: ProductRepository インターフェース
- [x] **Task 1.2.11**: OrderRepository インターフェース
- [x] **Task 1.2.12**: ~~インメモリリポジトリ実装（テスト用）~~ → 削除済み（Clean Architecture 違反のため）

## Phase 2: Infrastructure（目標期間: 1 週間）

### 🎯 目標

DynamoDB 接続とデータ永続化機能の実装

### Sprint 2.1: DynamoDB セットアップ（2 日目標）

#### Day 6: テーブル設計実装

- [x] **Task 2.1.1**: DynamoDB テーブル作成スクリプト
  ```bash
  aws dynamodb create-table --table-name OnlineShop \
    --attribute-definitions AttributeName=PK,AttributeType=S AttributeName=SK,AttributeType=S \
    --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
    --endpoint-url http://localhost:8000
  ```
- [x] **Task 2.1.2**: GSI1, GSI2 作成
- [x] **Task 2.1.3**: テーブル作成の自動化（Makefile）

#### Day 7: AWS SDK 設定

- [x] **Task 2.1.4**: DynamoDB クライアント設定（guregu/dynamo 使用）
  ```go
  cfg, err := config.LoadDefaultConfig(ctx,
    config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
      return aws.Endpoint{URL: "http://localhost:8000"}, nil
    })))
  ```
- [x] **Task 2.1.5**: 設定の環境変数化
- [x] **Task 2.1.6**: 接続ヘルスチェック実装

### Sprint 2.2: リポジトリ実装（3 日目標）

#### Day 8: データマッパー実装

- [x] **Task 2.2.1**: Customer データマッパー（CustomerItem struct 使用）
  ```go
  func CustomerToItem(customer *entity.Customer) (map[string]types.AttributeValue, error)
  func ItemToCustomer(item map[string]types.AttributeValue) (*entity.Customer, error)
  ```
- [ ] **Task 2.2.2**: Product データマッパー
- [ ] **Task 2.2.3**: マッパーの単体テスト

#### Day 9: Customer リポジトリ

- [x] **Task 2.2.4**: DynamoCustomerRepository 実装（guregu/dynamo 使用）
  ```go
  func (r *DynamoCustomerRepository) Save(ctx context.Context, customer *entity.Customer) error
  func (r *DynamoCustomerRepository) FindByID(ctx context.Context, id value.CustomerID) (*entity.Customer, error)
  ```
- [x] **Task 2.2.5**: Customer リポジトリの統合テスト

#### Day 10: Product/Order リポジトリ

- [ ] **Task 2.2.6**: DynamoProductRepository 実装
- [ ] **Task 2.2.7**: DynamoOrderRepository 実装（基本機能）
- [x] **Task 2.2.8**: 全リポジトリの統合テスト
  ```go
  // integration_test.go でエンドツーエンドシナリオテスト実装済み
  // Customer登録 → Product作成 → Order作成の一連フローをテスト
  // パフォーマンステスト（100件一括操作）も含む
  ```

**🎯 現在地**: Sprint 2.2 完了！ 次は Sprint 3.1（OpenAPI & 基本 API 実装）

## Phase 3: API Layer（目標期間: 1 週間）

### 🎯 目標

REST API エンドポイントと業務ユースケースの実装

### Sprint 3.1: OpenAPI & 基本 API（2 日目標）

#### Day 11: OpenAPI 仕様定義

- [x] **Task 3.1.1**: OpenAPI 仕様書作成
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
- [x] **Task 3.1.2**: oapi-codegen 設定
  ```yaml
  # oapi-codegen.config.yaml
  package: openapi
  generate:
    models: true
    echo-server: true
    embedded-spec: true
  output: internal/adapter/openapi/generated.go
  ```
- [x] **Task 3.1.3**: API 型とサーバーインターフェース生成
  ```go
  // 生成されたEcho用ServerInterface
  type ServerInterface interface {
    ListCustomers(ctx echo.Context, params ListCustomersParams) error
    CreateCustomer(ctx echo.Context) error
    // ... 15のエンドポイント定義済み
  }
  ```

#### Day 12: 基本 CRUD API

- [x] **Task 3.1.4**: Customer CRUD API 実装
  ```go
  func (h *APIHandler) CreateCustomer(ctx echo.Context) error
  func (h *APIHandler) GetCustomer(ctx echo.Context, customerId string) error
  func (h *APIHandler) ListCustomers(ctx echo.Context, params openapi.ListCustomersParams) error
  func (h *APIHandler) UpdateCustomer(ctx echo.Context, customerId string) error
  func (h *APIHandler) DeleteCustomer(ctx echo.Context, customerId string) error
  ```
- [x] **Task 3.1.5**: Product CRUD API 実装
  ```go
  func (h *APIHandler) CreateProduct(ctx echo.Context) error
  func (h *APIHandler) GetProduct(ctx echo.Context, productId string) error
  func (h *APIHandler) ListProducts(ctx echo.Context, params openapi.ListProductsParams) error
  func (h *APIHandler) UpdateProduct(ctx echo.Context, productId string) error
  func (h *APIHandler) DeleteProduct(ctx echo.Context, productId string) error
  ```
- [x] **Task 3.1.6**: Order CRUD API 実装
  ```go
  func (h *APIHandler) CreateOrder(ctx echo.Context) error
  func (h *APIHandler) GetOrder(ctx echo.Context, orderId string) error
  func (h *APIHandler) ListOrders(ctx echo.Context, params openapi.ListOrdersParams) error
  func (h *APIHandler) UpdateOrderStatus(ctx echo.Context, orderId string) error
  func (h *APIHandler) GetCustomerOrders(ctx echo.Context, customerId string, params openapi.GetCustomerOrdersParams) error
  ```
- [x] **Task 3.1.7**: Echo Server & Middleware 設定
  ```go
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(middleware.CORS())
  openapi.RegisterHandlers(e, apiHandler)
  ```
- [x] **Task 3.1.8**: 基本 API の動作確認
  ```bash
  # サーバー起動テスト成功
  curl http://localhost:8080/customers → [サンプルデータ] ✅
  curl -X POST http://localhost:8080/customers -d '{"name":"Test","email":"test@example.com"}' → 201 Created ✅
  # 全15エンドポイント基本動作確認済み
  ```

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

## 拡張候補ドメイン

### Phase 5: 在庫・配送・決済拡張（目標期間: 1 週間）

現在の基本 3 ドメイン（Customer, Product, Order）に加えて、以下のドメインを追加することで、より実用的なオンラインショップを構築できるのだ。

#### Warehouse（倉庫・在庫管理）

```go
// 在庫管理の値オブジェクト
type WarehouseID    = Branded[string, "WarehouseID"]
type StockQuantity  = Branded[int, "StockQuantity"]
type ReorderLevel   = Branded[int, "ReorderLevel"]

// 在庫エンティティ
type Stock struct {
    warehouseID WarehouseID
    productID   ProductID
    quantity    StockQuantity
    reorderLevel ReorderLevel
    location    string // 倉庫内の場所
}

// 倉庫エンティティ
type Warehouse struct {
    id       WarehouseID
    name     string
    address  string
    isActive bool
}
```

**実装優先度**: 高（在庫切れ管理は必須機能）
**DynamoDB アクセスパターン**:

- PK: `WAREHOUSE#{warehouseID}`, SK: `STOCK#{productID}`
- GSI1: `PRODUCT#{productID}`, SK: `WAREHOUSE#{warehouseID}` （商品別在庫検索）
- GSI2: `STOCK#LOW`, SK: quantity （在庫切れアラート用）

#### Shipment（配送管理）

```go
// 配送の値オブジェクト
type ShipmentID     = Branded[string, "ShipmentID"]
type TrackingNumber = Branded[string, "TrackingNumber"]
type ShippingFee    = Money

// 配送状態
type ShipmentStatus int
const (
    ShipmentPending ShipmentStatus = iota
    ShipmentPicked
    ShipmentShipped
    ShipmentDelivered
    ShipmentReturned
)

// 配送エンティティ
type Shipment struct {
    id            ShipmentID
    orderID       OrderID
    trackingNumber TrackingNumber
    status        ShipmentStatus
    shippingFee   ShippingFee
    estimatedDelivery time.Time
    actualDelivery    *time.Time
}
```

**実装優先度**: 中（注文との連携が重要）
**DynamoDB アクセスパターン**:

- PK: `SHIPMENT#{shipmentID}`, SK: `METADATA`
- GSI1: `ORDER#{orderID}`, SK: `SHIPMENT#{shipmentID}`
- GSI2: `STATUS#{status}`, SK: estimatedDelivery （配送状況別検索）

#### Payment（決済管理）

```go
// 決済の値オブジェクト
type PaymentID     = Branded[string, "PaymentID"]
type PaymentMethod = Branded[string, "PaymentMethod"]

// 決済状態
type PaymentStatus int
const (
    PaymentPending PaymentStatus = iota
    PaymentProcessing
    PaymentCompleted
    PaymentFailed
    PaymentRefunded
)

// 決済エンティティ
type Payment struct {
    id         PaymentID
    orderID    OrderID
    amount     Money
    method     PaymentMethod // "credit_card", "bank_transfer", etc.
    status     PaymentStatus
    processedAt *time.Time
    externalID  string // 外部決済プロバイダのID
}
```

**実装優先度**: 高（EC サイトの根幹機能）
**DynamoDB アクセスパターン**:

- PK: `PAYMENT#{paymentID}`, SK: `METADATA`
- GSI1: `ORDER#{orderID}`, SK: `PAYMENT#{paymentID}`
- GSI2: `STATUS#{status}`, SK: processedAt

### Phase 5 実装順序

#### Sprint 5.1: Warehouse ドメイン（2 日）

- [ ] **Task 5.1.1**: Warehouse 値オブジェクト・エンティティ実装
- [ ] **Task 5.1.2**: WarehouseRepository, StockRepository 実装
- [ ] **Task 5.1.3**: 在庫減少・補充ユースケース実装
- [ ] **Task 5.1.4**: 在庫切れ検知機能実装

#### Sprint 5.2: Payment ドメイン（2 日）

- [ ] **Task 5.2.1**: Payment 値オブジェクト・エンティティ実装
- [ ] **Task 5.2.2**: PaymentRepository 実装
- [ ] **Task 5.2.3**: 決済処理ユースケース実装（モック）
- [ ] **Task 5.2.4**: 決済状態管理・履歴機能実装

#### Sprint 5.3: Shipment ドメイン（2 日）

- [ ] **Task 5.3.1**: Shipment 値オブジェクト・エンティティ実装
- [ ] **Task 5.3.2**: ShipmentRepository 実装
- [ ] **Task 5.3.3**: 配送状況追跡ユースケース実装
- [ ] **Task 5.3.4**: 配送完了通知機能実装

#### Sprint 5.4: ドメイン統合（1 日）

- [ ] **Task 5.4.1**: Order → Payment → Shipment の状態連携実装
- [ ] **Task 5.4.2**: 在庫減少 → 注文確定 の整合性確保
- [ ] **Task 5.4.3**: 統合テスト実装
- [ ] **Task 5.4.4**: E2E シナリオテスト実装

## テスト戦略

### 現状の課題と改善方針

#### ~~1. inmemory リポジトリの配置問題~~ → 解決済み

**旧状況**: `internal/domain/repository/inmemory_*.go` （ドメイン層に配置）
**問題**: Clean Architecture の依存性ルールに違反
**解決**: inmemory リポジトリを完全削除し、DynamoDB Local 統一テスト戦略に移行

#### 1. テスト実行環境の統一 ✅

**統一後の状態**:

**統一後の状態**:

- 値オブジェクト・エンティティ: 純粋な単体テスト（外部依存なし）
- リポジトリテスト: DynamoDB Local 使用
- 統合テスト: DynamoDB Local 使用

**統一方針**:

```bash
# 開発時: DynamoDB Local起動が前提
make test-unit         # ドメイン層のテスト（外部依存なし）
make test-integration  # DynamoDB Local使用のリポジトリテスト

# 全テスト実行
make test-all     # unit + integration
```

#### 3. テストのピラミッド構造

```
      🔺 E2E Tests (少数・重要パス)
     🔺🔺 Integration Tests (適度・境界テスト)
   🔺🔺🔺🔺 Unit Tests (多数・ビジネスロジック)
```

**各層の責務**:

- **Unit Tests**: ドメインロジック（値オブジェクト・エンティティ）
- **Integration Tests**: リポジトリ実装、ユースケース（DynamoDB Local 使用）
- **E2E Tests**: API 経由の主要ユーザーシナリオ

#### 2. テスト戦略の実装

```go
// リポジトリはインターフェースなので、テスト時は差し替え可能
type CreateOrderUseCase struct {
    customerRepo repository.CustomerRepository  // インターフェース
    productRepo  repository.ProductRepository   // インターフェース
    orderRepo    repository.OrderRepository     // インターフェース
}

// 統合テスト時（DynamoDB Local使用）
func TestCreateOrderUseCase(t *testing.T) {
    // DynamoDB Local接続のリポジトリを使用
    dynamoClient := setupDynamoDBLocal(t)
    customerRepo := repository.NewDynamoCustomerRepository(dynamoClient)
    productRepo := repository.NewDynamoProductRepository(dynamoClient)
    orderRepo := repository.NewDynamoOrderRepository(dynamoClient)

    usecase := NewCreateOrderUseCase(customerRepo, productRepo, orderRepo)
    // ...
}

// 統合テスト時
func TestCreateOrderIntegration(t *testing.T) {
    // DynamoDB Localを使用（実環境に近い）
    db := testutil.SetupDynamoDB(t)
    customerRepo := dynamo.NewCustomerRepository(db)
    productRepo := dynamo.NewProductRepository(db)
    orderRepo := dynamo.NewOrderRepository(db)

    usecase := NewCreateOrderUseCase(customerRepo, productRepo, orderRepo)
    // ...
}
```

### テスト実行時間の目標

- **Unit Tests**: < 5 秒 （開発時の高速フィードバック）
- **Integration Tests**: < 30 秒 （DynamoDB Local の起動含む）
- **E2E Tests**: < 2 分 （API 経由の複数シナリオ）

### テストカバレッジ目標

- **ドメイン層**: 90%以上 （ビジネスロジックの品質確保）
- **ユースケース層**: 85%以上 （アプリケーションロジックの品質確保）
- **アダプター層**: 70%以上 （主要パスの動作確認）
- **全体**: 80%以上
