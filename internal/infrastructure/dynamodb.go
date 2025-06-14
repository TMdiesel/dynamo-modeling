package infrastructure

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/guregu/dynamo/v2"
)

// DynamoDBConfig holds the configuration for DynamoDB connection
type DynamoDBConfig struct {
	Region    string
	Endpoint  string // For local development
	TableName string
}

// DynamoDBClient wraps the guregu dynamo client
type DynamoDBClient struct {
	DB        *dynamo.DB
	TableName string
}

// NewDynamoDBClient creates a new DynamoDB client using guregu/dynamo
func NewDynamoDBClient(ctx context.Context, cfg DynamoDBConfig) (*DynamoDBClient, error) {
	slog.Info("Initializing DynamoDB client",
		"region", cfg.Region,
		"endpoint", cfg.Endpoint,
		"tableName", cfg.TableName)

	// AWS SDK config
	var awsCfg aws.Config
	var err error

	if cfg.Endpoint != "" {
		// For local development with DynamoDB Local
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           cfg.Endpoint,
						SigningRegion: cfg.Region,
					}, nil
				})),
		)
	} else {
		// For production with real AWS DynamoDB
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
		)
	}

	if err != nil {
		return nil, err
	}

	// Create DynamoDB service client
	dynamoSvc := dynamodb.NewFromConfig(awsCfg)

	// Create guregu dynamo DB client
	db := dynamo.NewFromIface(dynamoSvc)

	slog.Info("DynamoDB client initialized successfully")

	return &DynamoDBClient{
		DB:        db,
		TableName: cfg.TableName,
	}, nil
}

// GetTable returns a dynamo table instance
func (c *DynamoDBClient) GetTable() dynamo.Table {
	return c.DB.Table(c.TableName)
}

// HealthCheck performs a basic health check on the DynamoDB connection
func (c *DynamoDBClient) HealthCheck(ctx context.Context) error {
	// Try to describe the table to check connectivity
	_, err := c.GetTable().Describe().Run(ctx)
	if err != nil {
		slog.Error("DynamoDB health check failed", "error", err)
		return err
	}

	slog.Info("DynamoDB health check passed")
	return nil
}
