package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

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
