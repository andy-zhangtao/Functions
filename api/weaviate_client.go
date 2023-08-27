package handler

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

type weaviateClient struct {
	client *weaviate.Client
}

func NewWeaviateClient(host, schema, key string) (*weaviateClient, error) {
	cfg := weaviate.Config{
		Host:       "localhost:8080",
		Scheme:     "http",
		AuthConfig: auth.ApiKey{Value: key},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create weaviate client: %v", err)
	}

	return &weaviateClient{
		client: client,
	}, nil
}

func (wc *weaviateClient) AddNewRecord(class string, properties map[string]string) error {
	data := make(map[string]string)

	for key, val := range properties {
		data[key] = val
	}

	created, err := wc.client.Data().Creator().WithClassName(class).WithProperties(data).Do(context.Background())
	if err != nil {
		return fmt.Errorf("could not create record: %v", err)
	}

	logrus.Infof("Created record with id [%+v]", created)
	return nil
}
