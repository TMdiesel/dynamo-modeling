# Makefile for DynamoDB + Clean Architecture Project

.PHONY: help setup run test clean docker-up docker-down create-table generate

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  setup           - 開発環境のセットアップ"
	@echo "  generate        - OpenAPIからGoコードを生成"
	@echo "  run             - アプリケーションを起動"
	@echo "  test            - テストを実行"
	@echo "  test-coverage   - テストカバレッジを確認"
	@echo "  docker-up       - DynamoDB Local + Admin GUIをdocker-composeで起動"
	@echo "  docker-down     - docker-composeを停止"
	@echo "  admin           - DynamoDB Admin GUIをブラウザで開く"
	@echo "  test-connection - DynamoDB Local接続テスト"
	@echo "  create-table    - DynamoDBにテーブルを作成"
	@echo "  list-tables     - DynamoDBのテーブル一覧を表示"
	@echo "  clean           - ビルド成果物を削除"

# 開発環境セットアップ
setup:
	go mod tidy
	go mod download

# OpenAPIからGoコードを生成
generate:
	@echo "Generating OpenAPI code..."
	oapi-codegen --config=oapi-codegen.config.yaml api/openapi.yml

# アプリケーション起動
run:
	go run cmd/server/main.go

# テスト実行
test:
	go test ./...

# テストカバレッジ確認
test-coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# DynamoDB Local起動
docker-up:
	docker-compose up -d
	@echo "Waiting for DynamoDB Local to start..."
	@sleep 5

# Docker停止
docker-down:
	docker-compose down

# DynamoDB Admin GUIをブラウザで開く
admin:
	@echo "Opening DynamoDB Admin GUI..."
	@echo "DynamoDB Admin GUI: http://localhost:8001"
	@if command -v open >/dev/null 2>&1; then \
		open http://localhost:8001; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open http://localhost:8001; \
	else \
		echo "Please open http://localhost:8001 in your browser"; \
	fi

# DynamoDBテーブル作成
create-table:
	@echo "Creating OnlineShop table..."
	go run scripts/create_tables.go

# DynamoDB接続テスト
test-connection:
	@echo "Testing DynamoDB Local connection..."
	go run scripts/test_connection.go

# テーブル一覧確認
list-tables:
	@echo "Listing DynamoDB tables..."
	go run scripts/list_tables.go
	aws dynamodb list-tables --endpoint-url http://localhost:8000 --region us-east-1

# ビルド成果物削除
clean:
	go clean
	rm -f coverage.out coverage.html

# 依存関係更新
deps-update:
	go get -u ./...
	go mod tidy
