//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	slog.Info("DynamoDB Local接続テストを開始します")

	// DynamoDB Local用のクライアント設定
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localhost:8000",
					SigningRegion: "ap-northeast-1",
				}, nil
			})),
	)
	if err != nil {
		log.Fatalf("AWS設定の読み込みに失敗: %v", err)
	}

	// DynamoDBクライアントを作成
	client := dynamodb.NewFromConfig(cfg)

	// 接続テスト: テーブル一覧を取得
	slog.Info("DynamoDB Localに接続中...")
	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("DynamoDB Localへの接続に失敗: %v", err)
	}

	slog.Info("接続成功！現在のテーブル一覧:", "tables", result.TableNames)

	// テストテーブルを作成してみる
	testTableName := "ConnectionTest"
	slog.Info("テストテーブルを作成中...", "tableName", testTableName)

	_, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		// テーブルが既に存在する場合はエラーにならない
		slog.Warn("テーブル作成でエラーが発生", "error", err)
	} else {
		slog.Info("テストテーブルの作成に成功", "tableName", testTableName)
	}

	// テーブル一覧を再度取得
	result, err = client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("テーブル一覧の取得に失敗: %v", err)
	}

	slog.Info("最終的なテーブル一覧:", "tables", result.TableNames)
	fmt.Println("✅ DynamoDB Local接続テストが正常に完了しました！")
}
