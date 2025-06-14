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
)

func main() {
	slog.Info("DynamoDBテーブル一覧を取得します")

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

	// テーブル一覧を取得
	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("テーブル一覧の取得に失敗: %v", err)
	}

	if len(result.TableNames) == 0 {
		fmt.Println("📄 テーブルが見つかりませんでした")
		return
	}

	fmt.Printf("📋 DynamoDBテーブル一覧 (%d個)\n", len(result.TableNames))
	fmt.Println("=" + fmt.Sprintf("%50s", "="))

	for i, tableName := range result.TableNames {
		fmt.Printf("%d. %s\n", i+1, tableName)

		// テーブルの詳細情報を取得
		desc, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			slog.Warn("テーブル詳細の取得に失敗", "tableName", tableName, "error", err)
			continue
		}

		table := desc.Table
		fmt.Printf("   ステータス: %s\n", table.TableStatus)
		fmt.Printf("   アイテム数: %d\n", *table.ItemCount)
		fmt.Printf("   サイズ: %d bytes\n", *table.TableSizeBytes)

		// キースキーマ
		fmt.Printf("   キー: ")
		for j, key := range table.KeySchema {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s (%s)", *key.AttributeName, key.KeyType)
		}
		fmt.Println()

		// GSI
		if len(table.GlobalSecondaryIndexes) > 0 {
			fmt.Printf("   GSI: %d個\n", len(table.GlobalSecondaryIndexes))
			for _, gsi := range table.GlobalSecondaryIndexes {
				fmt.Printf("     - %s: ", *gsi.IndexName)
				for k, key := range gsi.KeySchema {
					if k > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s (%s)", *key.AttributeName, key.KeyType)
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}
