# Active Context

## 現在のフォーカス

Sprint 1.1（環境構築）が完了し、次に Sprint 1.2（ドメイン層実装）を開始する準備が整いました。

## 直近の作業内容

1. **Phase 1.1: 環境構築完了**

   - Go module 初期化とプロジェクト構造作成
   - Docker Compose で DynamoDB Local 環境構築
   - テーブル作成スクリプト作成・実行
   - 接続テスト・テーブル一覧確認スクリプト完成

2. **技術的成果**

   - OnlineShop テーブル作成（Single Table Design）
   - PK/SK + GSI1/GSI2 の構成完了
   - 開発用 Makefile 完成（自動化コマンド）
   - 基本依存関係設定（Echo, AWS SDK v2, testify）

3. **学習事項**
   - AWS SDK v2 での DynamoDB Local 接続方法
   - GSI 設定の制約事項（BillingMode 指定不可）
   - waiter 機能の正しい使用方法

## 次のステップ

1. **Sprint 1.2: ドメイン層実装開始**

   - Task 1.2.1: 値オブジェクトの実装（Money, Email, ProductId 等）
   - Task 1.2.2: エンティティの実装（Product, Customer, Order 等）
   - Task 1.2.3: ドメインエラーの定義
   - Task 1.2.4: リポジトリインターフェースの定義

2. **実装の指針**

   - 関数型アプローチ（純粋関数・不変データ）
   - ドメイン駆動設計（値オブジェクト & エンティティ）
   - テスト駆動開発（Red → Green → Refactor）
   - 早期リターンでネストを浅く

3. **準備完了事項**
   - Clean Architecture ディレクトリ構造
   - DynamoDB Local 動作環境
   - テストフレームワーク（testify）
   - 開発ツール（Makefile、各種スクリプト）

## 重要な設計決定

1. **Single Table Design 採用**

   - 理由: DynamoDB のベストプラクティスに従い、パフォーマンス最適化
   - トレードオフ: 複雑性増加、但し学習目的には最適

2. **Branded Type 使用**

   - CustomerID, ProductID 等の型安全性確保
   - コンパイル時の型チェックでバグ防止

3. **値オブジェクトによるドメインルール表現**
   - Email, Money, CustomerName 等での自己検証
   - ドメインロジックのカプセル化

## 学習ポイント

- DynamoDB のアクセスパターン駆動設計の重要性
- Clean Architecture での依存性逆転の実現方法
- NoSQL とドメインモデルのマッピング戦略
- Mermaid 図表による設計の可視化効果
- Single Table Design の複雑性と設計パターン理解
