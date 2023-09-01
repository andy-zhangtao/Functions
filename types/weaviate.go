package types

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data"
)

const (
	RequestVersionDefault = RequestVersionV1
	RequestVersionV1      = "v1"
)

const (
	EnvWeaviateHost   = "WEAVIATE_HOST"
	EnvWeaviateSchema = "WEAVIATE_SCHEMA"
	EnvWewaviateKey   = "WEAVIATE_KEY"
)

const (
	DiaryClassName = "Diary"
)

type WeaviateClient struct {
	client *weaviate.Client
}

func NewWeaviateClient(host, schema, key string) (*WeaviateClient, error) {
	cfg := weaviate.Config{
		Host:       host,
		Scheme:     schema,
		AuthConfig: auth.ApiKey{Value: key},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create weaviate client: %v", err)
	}

	return &WeaviateClient{
		client: client,
	}, nil
}

func (wc *WeaviateClient) AddNewRecord(class string, properties map[string]string) (*data.ObjectWrapper, error) {
	data := make(map[string]string)

	for key, val := range properties {
		data[key] = val
	}

	created, err := wc.client.Data().Creator().WithClassName(class).WithProperties(data).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not create record: %v", err)
	}

	logrus.Infof("Created record with id [%+v]", created)
	return created, nil
}

type DirayCreateModel struct {
	User    string      `json:"user"`
	Body    string      `json:"body"`
	Date    string      `json:"date,omitempty"`
	Tags    []string    `json:"tags"`
	Mask    interface{} `json:"mask,omitempty"`
	Version string      `json:"version"`
}

type DirayCreateResponse struct {
	Version string `json:"version"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

type DirayQueryModel struct {
	Version string   `json:"version"`
	User    string   `json:"user"`
	Date    string   `json:"date,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Keys    []string `json:"keys,omitempty"`
}

type DirayQueryResponse struct {
	Version string   `json:"version"`
	Status  string   `json:"status"`
	Code    int      `json:"code"`
	Records []string `json:"records"`
}
