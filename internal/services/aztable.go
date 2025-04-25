package services

import (
	"context"

	"encoding/json"
	"errors"
	"fmt"
	"ws-middle-cache/internal/middleware"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type AzureTableClient struct {
	serviceClient *aztables.ServiceClient
	tableClient   *aztables.Client
	tableName     string
}

// NewAzureTableClient initializes a new AzureTableClient
func NewAzureTableClient(accountName, accountKey, tableName string) (*AzureTableClient, error) {
	logger := middleware.NewLogger()
	credential, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}

	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/", accountName)
	serviceClient, err := aztables.NewServiceClientWithSharedKey(serviceURL, credential, nil)
	if err != nil {
		return nil, err
	}

	tableClient := serviceClient.NewClient(tableName)

	// Ensure the table exists
	_, err = tableClient.CreateTable(context.Background(), nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 409 { // 409 indicates ResourceExists

			logger.Info("Table already exists:", tableName)
		}
	}

	return &AzureTableClient{
		serviceClient: serviceClient,
		tableClient:   tableClient,
		tableName:     tableName,
	}, nil
}

// InsertEntity inserts a new entity into the table
func (c *AzureTableClient) InsertEntity(ctx context.Context, partitionKey, rowKey string, properties map[string]any) error {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: partitionKey,
			RowKey:       rowKey,
		},
		Properties: properties,
	}

	entityBytes, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	_, err = c.tableClient.AddEntity(ctx, entityBytes, nil)
	return err
}

// GetEntity retrieves an entity from the table
func (c *AzureTableClient) GetEntity(ctx context.Context, partitionKey, rowKey string) (*aztables.EDMEntity, error) {
	resp, err := c.tableClient.GetEntity(ctx, partitionKey, rowKey, nil)
	if err != nil {
		return nil, err
	}

	var entity aztables.EDMEntity
	err = json.Unmarshal(resp.Value, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// DeleteEntity deletes an entity from the table
func (c *AzureTableClient) DeleteEntity(ctx context.Context, partitionKey, rowKey string) error {
	_, err := c.tableClient.DeleteEntity(ctx, partitionKey, rowKey, nil)
	return err
}
