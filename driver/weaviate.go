package driver

import (
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

type WeaviateClient struct {
	client *weaviate.Client
}

type WeaviateClientConf struct {
	Host   string
	Schema string
	Key    string
}

func NewWeaviateClient(conf WeaviateClientConf) (*WeaviateClient, error) {
	cfg := weaviate.Config{
		Host:       conf.Host,
		Scheme:     conf.Schema,
		AuthConfig: auth.ApiKey{Value: conf.Key},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create weaviate client: %v", err)
	}

	return &WeaviateClient{
		client: client,
	}, nil
}
