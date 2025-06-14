# 実装計画書

**目次は自動生成なので編集しない**

- [実装計画書](#実装計画書)
  - [概要](#概要)
  - [実装方針](#実装方針)
    - [基本方針](#基本方針)
    - [📊 Sprint 3.2 完了サマリー](#-sprint-32-完了サマリー)
    - [🚨 **エラーハンドリング・冪等性・リトライ課題分析**](#-エラーハンドリング冪等性リトライ課題分析)
      - [1. **エラーハンドリングの課題**](#1-エラーハンドリングの課題)
      - [2. **冪等性の課題**](#2-冪等性の課題)
      - [3. **リトライ安全性の課題**](#3-リトライ安全性の課題)
      - [4. **トランザクション境界の問題**](#4-トランザクション境界の問題)
      - [5. **具体的な改善プラン**](#5-具体的な改善プラン)
      - [6. **Phase 4 での改善項目**](#6-phase-4-での改善項目)
      - [7. **E2E テストで発見された具体的な問題**](#7-e2e-テストで発見された具体的な問題)
      - [8. **テスト戦略の拡張**](#8-テスト戦略の拡張)
      - [9. **優先度の高い修正項目（Phase 4 で必須）**](#9-優先度の高い修正項目phase-4-で必須)
      - [10. **実装の前提条件と制約**](#10-実装の前提条件と制約)
      - [12. **リスクアセスメント（Phase 4）**](#12-リスクアセスメントphase-4)
      - [13. **技術的負債の管理**](#13-技術的負債の管理)
      - [14. **品質ゲートの定義**](#14-品質ゲートの定義)
      - [11. **Phase 4 での技術的判断基準**](#11-phase-4-での技術的判断基準)
  - [Phase 1: Foundation（目標期間: 1 週間）](#phase-1-foundation目標期間-1-週間)
    - [🎯 目標](#-目標)
    - [Sprint 1.1: プロジェクト初期化（2 日目標）](#sprint-11-プロジェクト初期化2-日目標)
      - [Day 1: 環境セットアップ](#day-1-環境セットアップ)
      - [Day 2: DynamoDB 接続確認](#day-2-dynamodb-接続確認)
    - [Sprint 1.2: ドメイン層実装（3 日目標）](#sprint-12-ドメイン層実装3-日目標)
      - [Day 3: 値オブジェクト実装](#day-3-値オブジェクト実装)
      - [Day 4: エンティティ実装](#day-4-エンティティ実装)
      - [Day 5: リポジトリインターフェース](#day-5-リポジトリインターフェース)
  - [Phase 2: Infrastructure（目標期間: 1 週間）](#phase-2-infrastructure目標期間-1-週間)
    - [🎯 目標](#-目標-1)
    - [Sprint 2.1: DynamoDB セットアップ（2 日目標）](#sprint-21-dynamodb-セットアップ2-日目標)
      - [Day 6: テーブル設計実装](#day-6-テーブル設計実装)
      - [Day 7: AWS SDK 設定](#day-7-aws-sdk-設定)
    - [Sprint 2.2: リポジトリ実装（3 日目標）](#sprint-22-リポジトリ実装3-日目標)
      - [Day 8: データマッパー実装](#day-8-データマッパー実装)
      - [Day 9: Customer リポジトリ](#day-9-customer-リポジトリ)
      - [Day 10: Product/Order リポジトリ](#day-10-productorder-リポジトリ)
  - [Phase 3: API Layer（目標期間: 1 週間）](#phase-3-api-layer目標期間-1-週間)
    - [🎯 目標](#-目標-2)
    - [Sprint 3.1: OpenAPI \& 基本 API（2 日目標）](#sprint-31-openapi--基本-api2-日目標)
      - [Day 11: OpenAPI 仕様定義](#day-11-openapi-仕様定義)
      - [Day 12: 基本 CRUD API](#day-12-基本-crud-api)
    - [Sprint 3.2: 業務 API（3 日目標）](#sprint-32-業務-api3-日目標)
      - [Day 13: ユースケース層実装](#day-13-ユースケース層実装)
      - [Day 14: 注文機能実装](#day-14-注文機能実装)
      - [Day 15: エラーハンドリング・E2E テスト](#day-15-エラーハンドリングe2e-テスト)
  - [Phase 4: Quality \& Documentation（目標期間: 3 日）](#phase-4-quality--documentation目標期間-3-日)
    - [🎯 目標](#-目標-3)
      - [重要: 現在の優先順位](#重要-現在の優先順位)
    - [📊 Sprint 4.1: エラーハンドリング・冪等性改善（3 日）](#-sprint-41-エラーハンドリング冪等性改善3-日)
      - [Day 15: エラーハンドリング統一化](#day-15-エラーハンドリング統一化)
      - [Day 16: 冪等性・リトライ実装](#day-16-冪等性リトライ実装)
      - [Day 17: トランザクション・テスト強化](#day-17-トランザクションテスト強化)
    - [📊 Sprint 3.1 完了サマリー](#-sprint-31-完了サマリー)
    - [技術スタック確定](#技術スタック確定)
  - [進捗管理](#進捗管理)
    - [デイリーチェックポイント](#デイリーチェックポイント)
    - [ウィークリーレビュー](#ウィークリーレビュー)
    - [完了条件](#完了条件)
  - [緊急時対応](#緊急時対応)
    - [スケジュール遅延時](#スケジュール遅延時)
    - [技術的課題発生時](#技術的課題発生時)
    - [リソース不足時](#リソース不足時)
  - [拡張候補ドメイン](#拡張候補ドメイン)
    - [Phase 5: 在庫・配送・決済拡張（目標期間: 1 週間）](#phase-5-在庫配送決済拡張目標期間-1-週間)
      - [Warehouse（倉庫・在庫管理）](#warehouse倉庫在庫管理)
      - [Shipment（配送管理）](#shipment配送管理)
      - [Payment（決済管理）](#payment決済管理)
    - [Phase 5 実装順序](#phase-5-実装順序)
      - [Sprint 5.1: Warehouse ドメイン（2 日）](#sprint-51-warehouse-ドメイン2-日)
      - [Sprint 5.2: Payment ドメイン（2 日）](#sprint-52-payment-ドメイン2-日)
      - [Sprint 5.3: Shipment ドメイン（2 日）](#sprint-53-shipment-ドメイン2-日)
      - [Sprint 5.4: ドメイン統合（1 日）](#sprint-54-ドメイン統合1-日)
  - [テスト戦略](#テスト戦略)
    - [現状の課題と改善方針](#現状の課題と改善方針)
      - [~~1. inmemory リポジトリの配置問題~~ → 解決済み](#1-inmemory-リポジトリの配置問題--解決済み)
      - [1. テスト実行環境の統一 ✅](#1-テスト実行環境の統一-)
      - [3. テストのピラミッド構造](#3-テストのピラミッド構造)
      - [2. テスト戦略の実装](#2-テスト戦略の実装)
    - [テスト実行時間の目標](#テスト実行時間の目標)
    - [テストカバレッジ目標](#テストカバレッジ目標)

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

**🎯 現在地**: Sprint 3.2 完了！ Clean Architecture 実装が全ドメインで完了 🎊

### 📊 Sprint 3.2 完了サマリー

**Clean Architecture 実装完了**:

- ✅ Customer Domain: Handler → Controller → UseCase → Repository → Entity
- ✅ Product Domain: Handler → Controller → UseCase → Repository → Entity
- ✅ Order Domain: Handler → Controller → UseCase → Repository → Entity
- ✅ 全 Repository Interface 完全実装（FindByStatus, FindInStock 等）
- ✅ Entity State Management 修正（NewOrderWithState ファクトリー）
- ✅ Order Status 更新の永続化問題解決
- ✅ 在庫管理・注文ワークフロー実装
- ✅ 構造化ドメインエラーハンドリング
- ✅ 完全な依存性注入 (main.go)

**E2E テスト実装完了**:

- ✅ 実際の HTTP サーバーに対する End-to-End テスト実装
- ✅ 正常シナリオ: 顧客登録 → 商品作成 → 注文作成 → ステータス更新
- ✅ エラーシナリオ: 在庫不足、存在しない顧客での注文エラー
- ✅ 並行処理シナリオ: 複数の並行注文での在庫競合テスト
- ✅ 注文ライフサイクル: pending → confirmed → shipped → delivered

**API 動作確認**:

```bash
# 全ドメインCRUD動作確認済み
POST /customers → Customer作成 ✅
GET /customers/{id} → Customer取得 ✅
POST /products → Product作成 ✅
GET /products/{id} → Product取得 ✅
POST /orders → Order作成（在庫減少確認） ✅
PUT /orders/{id} → Order Status更新 ✅
# Order: pending → confirmed → shipped のワークフロー確認済み
```

**⚠️ 発見された課題**:

1. **在庫管理の競合制御問題**:

   - 並行注文テストで 3 つの注文が全て成功（本来は 1 つのみ成功すべき）
   - DynamoDB の ConditionExpression による楽観的ロックが未実装
   - 在庫数の整合性が保証されていない状態

2. **競合制御の実装が必要な箇所**:

   ```go
   // ProductRepository.UpdateStock で ConditionExpression 必要
   // 例: "stock >= :quantity" で在庫が十分な場合のみ更新

   // OrderRepository.Save でべき等性保証が必要
   // 例: OrderIDの重複チェック
   ```

3. **今後の改善方針**:
   - DynamoDB の Conditional Writes を活用した在庫管理実装
   - トランザクション処理の検討（DynamoDB Transactions）
   - 在庫確保 → 注文確定の 2 段階コミット実装

### 🚨 **エラーハンドリング・冪等性・リトライ課題分析**

#### 1. **エラーハンドリングの課題**

**現状の問題**:

- ドメインエラーは構造化されているが、インフラエラーとの混在が発生
- HTTP ステータスコードとドメインエラーの対応が不完全
- エラーレスポンスの統一性が不十分

**具体的な問題箇所**:

```go
// 現在の問題: エラー分類が曖昧
func (c *OrderController) CreateOrder(ctx echo.Context) error {
    order, err := c.createOrderUseCase.Execute(context.Background(), command)
    if err != nil {
        // 全てのエラーを400 BadRequestで返している
        return c.presenter.PresentError(ctx, http.StatusBadRequest, "creation_failed", err.Error())
    }
}

// 改善すべき点: エラータイプに応じた適切なHTTPステータス
// - CUSTOMER_NOT_FOUND → 404 Not Found
// - INSUFFICIENT_STOCK → 409 Conflict
// - REPOSITORY_ERROR → 500 Internal Server Error
```

#### 2. **冪等性の課題**

**現状の問題**:

- Order 作成が非冪等（同じリクエストで複数回実行すると重複注文が発生）
- Customer 作成で Email 重複チェックはあるが、競合時の処理が不完全
- Product 在庫更新が非冪等（同時更新で整合性が崩れる）

**冪等性が必要な操作**:

```go
// Order作成の冪等性実装が必要
type CreateOrderUseCase struct {
    // Idempotency Keyを受け取る仕組みが必要
}

// 実装すべき冪等性パターン:
// 1. クライアント指定のIdempotency Key
// 2. リクエストハッシュベースの重複検知
// 3. 操作結果のキャッシュ（24時間）
```

#### 3. **リトライ安全性の課題**

**現状の問題**:

- DynamoDB 接続エラー時のリトライが危険（在庫の二重減算リスク）
- 部分的な成功状態からの復旧処理がない
- タイムアウト処理が不十分

**リトライ可能 vs 不可能な操作**:

```go
// リトライ安全（冪等）
func (r *Repository) FindByID(ctx context.Context, id ID) error {
    // 読み取り操作：何度実行しても同じ結果
}

// リトライ危険（非冪等）
func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd Command) error {
    // 1. 在庫確認 ← リトライ安全
    // 2. 在庫減算 ← リトライ危険（ConditionExpression未実装）
    // 3. 注文保存 ← リトライ危険（重複チェック未実装）
}
```

#### 4. **トランザクション境界の問題**

**現状の問題**:

- 注文作成で「在庫減算」→「注文保存」が分離されており、中間失敗状態が発生可能
- ロールバック処理が未実装
- 部分的成功からの復旧ロジックがない

**改善が必要なトランザクション**:

```go
// 問題のあるトランザクション境界
func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error) {
    // ここで失敗すると在庫だけ減って注文が保存されない
    for _, item := range cmd.Items {
        product.ReserveStock(item.Quantity) // ←【危険】
        uc.productRepo.Save(ctx, product)   // ←【危険】
    }

    // ここで失敗すると在庫は減ったが注文が作成されない状態
    err := uc.orderRepo.Save(ctx, order)   // ←【危険】
    if err != nil {
        // ロールバック処理が未実装 ←【問題】
        return nil, err
    }
}
```

#### 5. **具体的な改善プラン**

**Phase 4.1: エラーハンドリング改善（1 日）**

```go
// 1. 構造化エラーレスポンス
type APIError struct {
    Code       string            `json:"code"`
    Message    string            `json:"message"`
    Details    map[string]string `json:"details,omitempty"`
    RetryAfter *int             `json:"retry_after,omitempty"`
}

// 2. エラータイプ別HTTPステータス
func (c *Controller) mapDomainErrorToHTTP(err error) (int, APIError) {
    switch {
    case isDomainError(err, "NOT_FOUND"):
        return http.StatusNotFound, APIError{...}
    case isDomainError(err, "INSUFFICIENT_STOCK"):
        return http.StatusConflict, APIError{...}
    case isDomainError(err, "VALIDATION_ERROR"):
        return http.StatusBadRequest, APIError{...}
    default:
        return http.StatusInternalServerError, APIError{...}
    }
}
```

**Phase 4.2: 冪等性実装（1 日）**

```go
// 1. Idempotency Key機能
type IdempotencyService interface {
    IsProcessed(ctx context.Context, key string) (bool, interface{}, error)
    StoreResult(ctx context.Context, key string, result interface{}) error
}

// 2. 注文作成の冪等性
func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error) {
    // クライアント指定またはリクエストハッシュベースのキー
    idempotencyKey := cmd.IdempotencyKey
    if idempotencyKey == "" {
        idempotencyKey = generateHash(cmd)
    }

    // 既に処理済みかチェック
    if processed, result, _ := uc.idempotency.IsProcessed(ctx, idempotencyKey); processed {
        return result.(*entity.Order), nil
    }

    // 実際の処理...
}
```

**Phase 4.3: 競合制御実装（1 日）**

```go
// 1. DynamoDB Conditional Writes
func (r *DynamoProductRepository) Save(ctx context.Context, product *entity.Product) error {
    item := ProductToItem(product)

    // 楽観的ロック：現在の在庫数が期待値と一致する場合のみ更新
    condition := expression.Name("Stock").Equal(expression.Value(product.PreviousStock()))

    err := r.table.Put(item).If(condition).Run(ctx)
    if err != nil {
        if isConditionalCheckFailedException(err) {
            return domain.NewDomainError("CONCURRENT_MODIFICATION",
                "Product was modified by another request", err)
        }
        return err
    }
    return nil
}

// 2. DynamoDB Transactions使用
func (uc *CreateOrderUseCase) ExecuteWithTransaction(ctx context.Context, cmd CreateOrderCommand) error {
    // すべての操作を単一トランザクションで実行
    transactionItems := []types.TransactWriteItem{
        // 在庫減算（条件付き）
        {Update: &types.Update{...}},
        // 注文保存
        {Put: &types.Put{...}},
    }

    _, err := r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
        TransactItems: transactionItems,
    })
    return err
}
```

#### 6. **Phase 4 での改善項目**

- [ ] **Task 4.1.1**: 構造化エラーレスポンス実装
- [ ] **Task 4.1.2**: エラータイプ別 HTTP ステータスマッピング
- [ ] **Task 4.1.3**: リトライ可能エラーの識別機能
- [ ] **Task 4.2.1**: Idempotency Key 機能実装
- [ ] **Task 4.2.2**: 注文作成の冪等性実装
- [ ] **Task 4.2.3**: 重複処理検知・結果キャッシュ
- [ ] **Task 4.3.1**: DynamoDB Conditional Writes 実装
- [ ] **Task 4.3.2**: 在庫管理の競合制御実装
- [ ] **Task 4.3.3**: DynamoDB Transactions 適用
- [ ] **Task 4.4.1**: サーキットブレーカー実装（基本）
- [ ] **Task 4.4.2**: エクスポネンシャルバックオフリトライ
- [ ] **Task 4.4.3**: 部分失敗からの復旧ロジック

#### 7. **E2E テストで発見された具体的な問題**

**並行処理での在庫管理問題**:

E2E テストの `TestE2E_ConcurrentOrders` で以下の問題が発覚:

- 在庫 5 個の商品に対して、3 つの並行注文（各 3 個）が全て成功
- 理想的には 1 つのみ成功し、残り 2 つは在庫不足エラーとなるべき
- 在庫が負数（-4 個）になる競合状態が発生

```go
// 実際の問題コード（internal/usecase/order_usecase.go）
func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error) {
    // 在庫確認（時点A）
    product, _ := uc.productRepo.FindByID(ctx, productID)

    // 在庫予約（時点B）- この間に他のリクエストが割り込み可能
    err = product.ReserveStock(itemCmd.Quantity)
    if err != nil {
        return nil, domain.NewDomainError("STOCK_RESERVATION_FAILED", ...)
    }

    // 在庫更新（時点C）- 楽観的ロックなしで更新
    err = uc.productRepo.Save(ctx, product)
    // ←【危険】: 他のリクエストの変更を上書きする可能性
}
```

**HTTP ステータスコードの不統一問題**:

```go
// 現在の問題例
func (c *OrderController) CreateOrder(ctx echo.Context) error {
    order, err := c.createOrderUseCase.Execute(context.Background(), command)
    if err != nil {
        // 全てのエラーが 400 BadRequest として返される
        return c.presenter.PresentError(ctx, http.StatusBadRequest, "creation_failed", err.Error())
    }
}

// 改善が必要:
// - CUSTOMER_NOT_FOUND → 404 Not Found
// - INSUFFICIENT_STOCK → 409 Conflict
// - CONCURRENT_MODIFICATION → 409 Conflict
// - VALIDATION_ERROR → 400 Bad Request
// - REPOSITORY_ERROR → 500 Internal Server Error
```

#### 8. **テスト戦略の拡張**

```go
// エラーハンドリングの統合テスト
func TestErrorHandling_ConcurrentStockUpdate(t *testing.T) {
    // 同時在庫更新でConflictエラーが適切に処理されることをテスト
}

func TestIdempotency_DuplicateOrderCreation(t *testing.T) {
    // 同じIdempotency Keyで複数回リクエストしても同じ結果が返ることをテスト
}

func TestRetry_PartialFailureRecovery(t *testing.T) {
    // 部分失敗状態からの復旧処理をテスト
}

func TestConcurrentStock_RaceCondition(t *testing.T) {
    // E2Eテストで発見された在庫競合問題の単体テスト版
    // 在庫5個に対して並行注文3つ（各3個）で適切にエラーハンドリングされることを確認
}
```

#### 9. **優先度の高い修正項目（Phase 4 で必須）**

**🚨 Critical レベル（必須修正）**:

1. **在庫管理の競合制御**

   - DynamoDB Conditional Writes による楽観的ロック実装
   - 理由: E2E テストで在庫がマイナスになる問題を発見

2. **構造化エラーレスポンス**

   - エラータイプ別 HTTP ステータスコードマッピング
   - 理由: 全エラーが 400 で返される問題

3. **Order 作成の冪等性**
   - Idempotency Key 機能実装
   - 理由: 重複注文防止が未実装

**⚠️ High レベル（推奨修正）**:

4. **DynamoDB Transactions 適用**

   - 在庫減算と注文保存の原子性保証
   - 理由: 部分失敗状態の回避

5. **リトライ機構の安全化**
   - 冪等性チェック付きリトライ実装
   - 理由: 危険な操作の重複実行防止

**📝 Medium レベル（時間があれば）**:

6. **サーキットブレーカー実装**
   - DynamoDB 接続エラー時の障害拡散防止
   - 理由: システム安定性向上

#### 10. **実装の前提条件と制約**

**技術的制約**:

- DynamoDB Local 環境での開発（本番は AWS DynamoDB）
- Single Table Design パターンの維持
- Clean Architecture 原則の遵守
- Go 1.21+の機能活用

**パフォーマンス要件**:

- API 応答時間: 95 パーセンタイル < 500ms
- 並行注文処理: 80%以上の成功率維持
- E2E テスト実行時間: < 30 秒

**可観測性要件**:

- 構造化ログ（JSON 形式）
- エラー種別の詳細情報
- パフォーマンスメトリクス取得可能性

#### 12. **リスクアセスメント（Phase 4）**

**🚨 高リスク（即座に対応が必要）**:

| リスク項目             | 影響度 | 発生確率 | 対策                             | 期限   |
| ---------------------- | ------ | -------- | -------------------------------- | ------ |
| 在庫競合問題           | 高     | 100%     | DynamoDB Conditional Writes 実装 | Day 15 |
| 重複注文               | 高     | 80%      | Idempotency Key 実装             | Day 16 |
| エラーレスポンス不統一 | 中     | 100%     | 構造化エラーレスポンス           | Day 15 |

**⚠️ 中リスク（計画的に対応）**:

| リスク項目         | 影響度 | 発生確率 | 対策                       | 期限     |
| ------------------ | ------ | -------- | -------------------------- | -------- |
| 部分失敗状態       | 高     | 40%      | DynamoDB Transactions 適用 | Day 17   |
| DynamoDB 接続障害  | 中     | 30%      | サーキットブレーカー実装   | Phase 後 |
| パフォーマンス劣化 | 中     | 50%      | ベンチマークテスト強化     | Day 17   |

**📝 低リスク（監視のみ）**:

- テストデータのクリーンアップ忘れ → テストスイート改善で対応
- 開発環境固有の問題 → 本番環境テストで検証
- ドキュメント更新漏れ → Phase 4 完了時に一括確認

#### 13. **技術的負債の管理**

**現在の技術的負債**:

1. **エラーハンドリングの一貫性不足**

   - 負債レベル: 高
   - 解決タイミング: Phase 4 必須
   - 影響: API の信頼性、デバッグ困難性

2. **冪等性の未実装**

   - 負債レベル: 高
   - 解決タイミング: Phase 4 必須
   - 影響: 重複処理、データ整合性

3. **ログ出力の構造化不足**

   - 負債レベル: 中
   - 解決タイミング: Phase 4 で部分改善
   - 影響: 運用監視の困難性

4. **テストカバレッジの偏り**
   - 負債レベル: 中
   - 解決タイミング: 継続的改善
   - 影響: バグ検出能力

**許容可能な負債**:

- パフォーマンス最適化（MVP 段階では不要）
- 詳細なメトリクス収集（基本機能安定後）
- 複雑なビジネスロジック（要件が固まってから）

#### 14. **品質ゲートの定義**

**Phase 4 完了の必須条件**:

✅ **機能品質**:

- [ ] 全 E2E テストが pass する
- [ ] 並行処理テストでエラーハンドリングが適切
- [ ] 冪等性テストが成功する
- [ ] レスポンス時間が要件を満たす

✅ **コード品質**:

- [ ] golint、go vet、gofmt が pass
- [ ] テストカバレッジ > 80%
- [ ] サイクロマティック複雑度 < 10
- [ ] 構造化ログが出力される

✅ **ドキュメント品質**:

- [ ] API 仕様書が最新
- [ ] README.md が更新済み
- [ ] エラーハンドリングガイド作成
- [ ] 運用手順書作成

#### 11. **Phase 4 での技術的判断基準**

**実装を進める基準**:

- 既存の E2E テストが通ること
- 新機能追加時は必ずテストも追加
- エラーハンドリングの改善が先行

**実装を後回しにする基準**:

- パフォーマンス最適化（機能完成後）
- 詳細なメトリクス収集（基本機能安定後）
- UI/UX 改善（API 安定後）

---

---

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

- [x] **Task 3.2.1**: CreateCustomerUseCase 実装
  ```go
  type CreateCustomerUseCase struct {
    customerRepo repository.CustomerRepository
  }
  func (uc *CreateCustomerUseCase) Execute(ctx context.Context, cmd CreateCustomerCommand) (*Customer, error)
  ```
- [x] **Task 3.2.2**: CreateProductUseCase 実装
- [x] **Task 3.2.3**: ユースケースの単体テスト

#### Day 14: 注文機能実装

- [x] **Task 3.2.4**: CreateOrderUseCase 実装
- [x] **Task 3.2.5**: GetCustomerOrdersUseCase 実装（GSI2 使用）
- [x] **Task 3.2.6**: 注文 API エンドポイント実装

#### Day 15: エラーハンドリング・E2E テスト

- [x] **Task 3.2.7**: 統一エラーハンドリング実装
- [x] **Task 3.2.8**: バリデーション実装
- [x] **Task 3.2.9**: E2E テスト実装

## Phase 4: Quality & Documentation（目標期間: 3 日）

### 🎯 目標

E2E テストで発見された**クリティカルな品質課題**の解決と基本的なドキュメント整備

#### 重要: 現在の優先順位

1. **エラーハンドリング・冪等性・競合制御** (Critical)
2. **基本的な品質向上・ドキュメント** (Important)

### 📊 Sprint 4.1: エラーハンドリング・冪等性改善（3 日）

#### Day 15: エラーハンドリング統一化

**重要**: E2E テストで発見された問題への対策を最優先で実装

- [ ] **Task 4.1.1**: 構造化エラーレスポンス実装

  - APIError 構造体定義（Code, Message, Details, RetryAfter）
  - ドメインエラーから HTTP ステータスへのマッピング関数
  - Presenter 層でのエラーレスポンス統一化

- [ ] **Task 4.1.2**: 在庫管理競合制御実装

  - DynamoDB Conditional Writes 実装
  - product.ReserveStock()での楽観的ロック追加
  - ConditionalCheckFailedException 処理

- [ ] **Task 4.1.3**: エラーハンドリングの単体テスト
  - 並行在庫更新での Conflict エラーテスト
  - HTTP ステータスマッピングテスト
  - エラーレスポンス形式テスト

#### Day 16: 冪等性・リトライ実装

- [ ] **Task 4.2.1**: Idempotency Key 機能実装

  - IdempotencyService interface 定義
  - DynamoDB ベースのキー管理実装
  - 24 時間 TTL 設定

- [ ] **Task 4.2.2**: Order 作成の冪等性実装

  - CreateOrderCommand に IdempotencyKey 追加
  - リクエストハッシュベースのキー生成
  - 処理済みリクエストの結果返却

- [ ] **Task 4.2.3**: リトライ安全性の改善
  - 冪等性チェック付きリトライ実装
  - エクスポネンシャルバックオフ
  - リトライ可能エラーの分類

#### Day 17: トランザクション・テスト強化

- [ ] **Task 4.3.1**: DynamoDB Transactions 適用

  - TransactWriteItems 使用での注文作成
  - 在庫減算と注文保存の原子性保証
  - トランザクション失敗時のエラーハンドリング

- [ ] **Task 4.3.2**: E2E テスト拡張

  - 冪等性テスト追加
  - 並行処理での整合性テスト強化
  - エラーレスポンス検証テスト

- [ ] **Task 4.3.3**: パフォーマンス指標測定
  - レスポンス時間の 95 パーセンタイル測定
  - 並行処理成功率の測定
  - リソース使用量の監視

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
