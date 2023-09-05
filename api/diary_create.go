package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	fmongo "github.com/andy-zhangtao/Functions/service/f_mongo"
	fweaviate "github.com/andy-zhangtao/Functions/service/f_weaviate"
	"github.com/andy-zhangtao/Functions/tools/flogs"
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
func DirayCreate(data string) (dcm types.DirayCreateModel, object *data.ObjectWrapper, err error) {

	// var dcm types.DirayCreateModel
	err = json.Unmarshal([]byte(data), &dcm)
	if err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		return dcm, object, err
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
		return dcm, object, fmt.Errorf("error parsing request body: %v", err)
	}

	dcm.DateSave = t

	wc, err := fweaviate.NewWeaviateClient(os.Getenv(types.EnvWeaviateHost), os.Getenv(types.EnvWeaviateSchema), os.Getenv(types.EnvWewaviateKey))
	if err != nil {
		logrus.Errorf("Error creating weaviate client: %v", err)
		return dcm, object, fmt.Errorf("error creating weaviate client: %v", err)
	}

	switch dcm.Version {
	case types.RequestVersionV1:
		err = checkV1(dcm)
		if err != nil {
			logrus.Errorf("Error parsing request body: %v", err)
			return dcm, object, fmt.Errorf("error parsing request body: %v", err)
		}

		object, err = wc.AddNewRecord(types.DiaryClassName, map[string]interface{}{
			"user":    dcm.User,
			"content": dcm.Body,
			"date":    dcm.DateSave.Unix(),
		})
		if err != nil {
			logrus.Errorf("Error creating weaviate record: %v", err)
			return dcm, object, fmt.Errorf("error creating weaviate record: %v", err)
		}

		return dcm, object, nil
	default:
		logrus.Errorf("Not support version: %v", dcm.Version)
		return dcm, object, fmt.Errorf("not support version: %v", dcm.Version)
	}
}

// DirayCreateHandler handle the diary create request
// @Summary create a new diary
// First save the diary to weaviate, then save the diary to mongo
func DirayCreateHandler(w http.ResponseWriter, r *http.Request) {

	// chech the http method , only allow POST method
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		flogs.Errorf("Error parsing request body: %v", err)
		errorResponse(w, err)
		return
	}

	flogs.Infof("request body: %v", string(data))

	dcm, object, err := DirayCreate(string(data))
	if err != nil {
		flogs.Errorf("Error parsing request body: %v", err)
		errorResponse(w, err)
		return
	}

	dcm.Mask = map[string]interface{}{
		"weaviate": object.Object.ID.String(),
	}

	err = createMongoRecord(dcm)
	if err != nil {
		flogs.Errorf("Error creating mongo record: %v", err)
		errorResponse(w, err)
		return
	}
	commonResponse(w, http.StatusOK, types.DirayCreateResponse{
		Code: http.StatusOK,
		Msg:  object.Object.ID.String(),
	})
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

// func createWeaviateRecord(dcm types.DirayCreateModel) error {
// 	return nil
// }

func createMongoRecord(dcm types.DirayCreateModel) error {
	flogs.Infof("createMongoRecord: %v", dcm)
	uri := os.Getenv(types.EnvMONGOHOST)
	db := os.Getenv(types.EnvMONGODB)
	collection := os.Getenv(types.EnvMONGOCOLLECTION)

	flogs.Infof("uri: %s db: %s collection: %s", uri, db, collection)
	cli, err := fmongo.NewMongoCli(uri, db, collection)
	if err != nil {
		return fmt.Errorf("create mongo client error: %w", err)
	}

	return cli.SaveDataToMongo(dcm, map[string]interface{}{})
}
