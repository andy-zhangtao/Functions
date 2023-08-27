package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
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

type weaviateClient struct {
	client *weaviate.Client
}

func NewWeaviateClient(host, schema, key string) (*weaviateClient, error) {
	cfg := weaviate.Config{
		Host:       host,
		Scheme:     schema,
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

type DirayCreateModel struct {
	User    string      `json:"user"`
	Body    string      `json:"body"`
	Tags    []string    `json:"tags"`
	Mask    interface{} `json:"mask,omitempty"`
	Version string      `json:"version"`
}

type DirayCreateResponse struct {
	Version string `json:"version"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

func DirayCreate(w http.ResponseWriter, r *http.Request) {

	// chech the http method , only allow POST method
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// get the body of the request
	var dcm DirayCreateModel
	err := json.NewDecoder(r.Body).Decode(&dcm)
	if err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		errorResponse(w, err)
		return
	}

	if dcm.Version == "" {
		dcm.Version = RequestVersionDefault
	}

	// output all envoriment variables
	for _, e := range os.Environ() {
		logrus.Infof("%v", e)
	}

	wc, err := NewWeaviateClient(os.Getenv(EnvWeaviateHost), os.Getenv(EnvWeaviateSchema), os.Getenv(EnvWewaviateKey))
	if err != nil {
		logrus.Errorf("Error creating weaviate client: %v", err)
		errorResponse(w, err)
		return
	}

	switch dcm.Version {
	case RequestVersionV1:
		err = checkV1(dcm)
		if err != nil {
			logrus.Errorf("Error parsing request body: %v", err)
			errorResponse(w, err)
			return
		}

		err = wc.AddNewRecord(DiaryClassName, map[string]string{
			"user":    dcm.User,
			"content": dcm.Body,
		})
		if err != nil {
			logrus.Errorf("Error creating weaviate record: %v", err)
			errorResponse(w, err)
			return
		}

		commonResponse(w, http.StatusOK, DirayCreateResponse{
			Code: http.StatusOK,
		})
	default:
		logrus.Errorf("Error parsing request body: %v", err)
		errorResponse(w, err)
		return
	}

}

// checkV1 check the request body for v1
func checkV1(dcm DirayCreateModel) error {
	if dcm.User == "" {
		return fmt.Errorf("user is empty")
	}

	if dcm.Body == "" {
		return fmt.Errorf("body is empty")
	}

	return nil
}

// errorResponse return the error response
func errorResponse(w http.ResponseWriter, err error) {
	commonResponse(w, http.StatusBadRequest, DirayCreateResponse{
		Code:    http.StatusBadRequest,
		Msg:     err.Error(),
		Version: RequestVersionDefault,
	})
}

// commonResponse
func commonResponse(w http.ResponseWriter, code int, data DirayCreateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(data)
}

func createWeaviateRecord(dcm DirayCreateModel) error {
	return nil
}
