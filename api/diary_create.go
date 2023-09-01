package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/andy-zhangtao/Functions/types"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data"
)

// DirayCreate create a new diary
// @Summary create a new diary
// @Description create a new diary
// @Tags diary
// @Accept  json
// @Produce  json
func DirayCreate(data string) (object *data.ObjectWrapper, err error) {

	var dcm types.DirayCreateModel
	err = json.Unmarshal([]byte(data), &dcm)
	if err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		return object, err
	}

	if dcm.Version == "" {
		dcm.Version = types.RequestVersionDefault
	}

	if dcm.Date == "" {
		// use yyyy-mm-dd format
		t := time.Now()
		dcm.Date = fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	}

	// parse dcm.Date to time.Time
	t, err := time.Parse("2006-01-02", dcm.Date)
	if err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		return object, fmt.Errorf("error parsing request body: %v", err)
	}

	dcm.DateSave = t

	wc, err := types.NewWeaviateClient(os.Getenv(types.EnvWeaviateHost), os.Getenv(types.EnvWeaviateSchema), os.Getenv(types.EnvWewaviateKey))
	if err != nil {
		logrus.Errorf("Error creating weaviate client: %v", err)
		return object, fmt.Errorf("error creating weaviate client: %v", err)
	}

	switch dcm.Version {
	case types.RequestVersionV1:
		err = checkV1(dcm)
		if err != nil {
			logrus.Errorf("Error parsing request body: %v", err)
			return object, fmt.Errorf("error parsing request body: %v", err)
		}

		object, err = wc.AddNewRecord(types.DiaryClassName, map[string]interface{}{
			"user":    dcm.User,
			"content": dcm.Body,
			"date":    dcm.DateSave,
		})
		if err != nil {
			logrus.Errorf("Error creating weaviate record: %v", err)
			return object, fmt.Errorf("error creating weaviate record: %v", err)
		}

		return object, nil
	default:
		logrus.Errorf("Not support version: %v", dcm.Version)
		return object, fmt.Errorf("not support version: %v", dcm.Version)
	}
}

func DirayCreateHandler(w http.ResponseWriter, r *http.Request) {

	// chech the http method , only allow POST method
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		errorResponse(w, err)
		return
	}

	logrus.Infof("request body: %v", string(data))

	object, err := DirayCreate(string(data))
	if err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		errorResponse(w, err)
		return
	}

	commonResponse(w, http.StatusOK, types.DirayCreateResponse{
		Code: http.StatusOK,
		Msg:  object.Object.ID.String(),
	})

	// get the body of the request
	// var dcm types.DirayCreateModel
	// err := json.NewDecoder(r.Body).Decode(&dcm)
	// if err != nil {
	// 	logrus.Errorf("Error parsing request body: %v", err)
	// 	errorResponse(w, err)
	// 	return
	// }

	// if dcm.Version == "" {
	// 	dcm.Version = types.RequestVersionDefault
	// }

	// // output all envoriment variables
	// for _, e := range os.Environ() {
	// 	logrus.Infof("%v", e)
	// }

	// wc, err := types.NewWeaviateClient(os.Getenv(types.EnvWeaviateHost), os.Getenv(types.EnvWeaviateSchema), os.Getenv(types.EnvWewaviateKey))
	// if err != nil {
	// 	logrus.Errorf("Error creating weaviate client: %v", err)
	// 	errorResponse(w, err)
	// 	return
	// }

	// switch dcm.Version {
	// case types.RequestVersionV1:
	// 	err = checkV1(dcm)
	// 	if err != nil {
	// 		logrus.Errorf("Error parsing request body: %v", err)
	// 		errorResponse(w, err)
	// 		return
	// 	}

	// 	err = wc.AddNewRecord(types.DiaryClassName, map[string]string{
	// 		"user":    dcm.User,
	// 		"content": dcm.Body,
	// 	})
	// 	if err != nil {
	// 		logrus.Errorf("Error creating weaviate record: %v", err)
	// 		errorResponse(w, err)
	// 		return
	// 	}

	// 	commonResponse(w, http.StatusOK, types.DirayCreateResponse{
	// 		Code: http.StatusOK,
	// 	})
	// default:
	// 	logrus.Errorf("Error parsing request body: %v", err)
	// 	errorResponse(w, err)
	// 	return
	// }

}

// checkV1 check the request body for v1
func checkV1(dcm types.DirayCreateModel) error {
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
	commonResponse(w, http.StatusBadRequest, types.DirayCreateResponse{
		Code:    http.StatusBadRequest,
		Msg:     err.Error(),
		Version: types.RequestVersionDefault,
	})
}

// commonResponse
func commonResponse(w http.ResponseWriter, code int, data types.DirayCreateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(data)
}

func createWeaviateRecord(dcm types.DirayCreateModel) error {
	return nil
}
