package orm

import (
    "log"
    "sync"

    "github.com/opensearch-project/opensearch-go/v2"
)

var (
    clientInstance *opensearch.Client
    once           sync.Once
)

// InitializeClient initializes the OpenSearch client
func InitializeClient(address string) error {
    var err error
    once.Do(func() {
        cfg := opensearch.Config{
            Addresses: []string{address},
        }
        clientInstance, err = opensearch.NewClient(cfg)
        if err != nil {
            log.Fatalf("Failed to create OpenSearch client: %v", err)
        }
    })
    return err
}

// GetClient retrieves the singleton OpenSearch client
func GetClient() *opensearch.Client {
    return clientInstance
}
