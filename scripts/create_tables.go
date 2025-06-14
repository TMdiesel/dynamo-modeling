//go:build ignore

// filepath: /Users/muratatakuya/project/aws/dynamo-modeling/scripts/create_tables.go
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	MainTableName = "OnlineShop"
	GSI1Name      = "GSI1"
	GSI2Name      = "GSI2"
)

func main() {
	slog.Info("オンラインショップテーブル作成スクリプトを開始します")

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

	client := dynamodb.NewFromConfig(cfg)

	// 既存テーブルをチェック
	existingTables, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("テーブル一覧の取得に失敗: %v", err)
	}

	// テーブルが既に存在するかチェック
	for _, tableName := range existingTables.TableNames {
		if tableName == MainTableName {
			slog.Warn("テーブルが既に存在します", "tableName", MainTableName)
			fmt.Printf("✅ テーブル '%s' は既に存在しています\n", MainTableName)
			return
		}
	}

	// メインテーブルを作成
	slog.Info("メインテーブルを作成中...", "tableName", MainTableName)

	createTableInput := &dynamodb.CreateTableInput{
		TableName: aws.String(MainTableName),

		// 属性定義
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI2PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI2SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},

		// キースキーマ（メインテーブル）
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},

		// GSI定義
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(GSI1Name),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("GSI1PK"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("GSI1SK"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
			{
				IndexName: aws.String(GSI2Name),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("GSI2PK"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("GSI2SK"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},

		// 従量課金制を使用
		BillingMode: types.BillingModePayPerRequest,
	}

	_, err = client.CreateTable(context.TODO(), createTableInput)
	if err != nil {
		log.Fatalf("テーブル作成に失敗: %v", err)
	}

	slog.Info("テーブル作成が開始されました", "tableName", MainTableName)

	// テーブルが作成されるまで待機
	waiter := dynamodb.NewTableExistsWaiter(client)
	err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(MainTableName),
	}, time.Minute*5) // 5分でタイムアウト
	if err != nil {
		log.Fatalf("テーブル作成の待機に失敗: %v", err)
	}

	// テーブルの詳細を取得して確認
	desc, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(MainTableName),
	})
	if err != nil {
		log.Fatalf("テーブル詳細の取得に失敗: %v", err)
	}

	slog.Info("テーブル作成完了",
		"tableName", *desc.Table.TableName,
		"status", desc.Table.TableStatus,
		"itemCount", desc.Table.ItemCount,
		"gsiCount", len(desc.Table.GlobalSecondaryIndexes),
	)

	fmt.Printf("✅ テーブル '%s' の作成が正常に完了しました！\n", MainTableName)
	fmt.Printf("   - メインキー: PK (Hash), SK (Range)\n")
	fmt.Printf("   - GSI1: GSI1PK (Hash), GSI1SK (Range)\n")
	fmt.Printf("   - GSI2: GSI2PK (Hash), GSI2SK (Range)\n")
	fmt.Printf("   - 課金モード: Pay per request\n")
}
