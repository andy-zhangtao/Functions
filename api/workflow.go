package handler

import (
	"net/http"
	"os"

	"github.com/andy-zhangtao/Functions/types"
	"github.com/andy-zhangtao/Functions/workflow"
)

// func WorkFlowHandler(w http.ResponseWriter, r *http.Request) {
// 	// chech the http method , only allow POST method
// 	if r.Method != "POST" {
// 		http.Error(w, "Method is not supported.", http.StatusNotFound)
// 		return
// 	}

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		flogs.Errorf("Error parsing request body: %v", err)
// 		errorResponse(w, err)
// 		return
// 	}

// 	flogs.Infof("request body: %v", string(data))

// 	var wfr types.WorkFlowRequest
// 	err = json.Unmarshal(data, &wfr)
// 	if err != nil {
// 		flogs.Errorf("Error parsing request body: %v", err)
// 		errorResponse(w, err)
// 		return
// 	}

// 	fwc, err := workflow()
// 	if err != nil {
// 		flogs.Errorf("Error creating workflow client: %v", err)
// 		errorResponse(w, err)
// 		return
// 	}

// 	switch wfr.Action {
// 	case types.WorkFlowActionGet:
// 		flows, err := fwc.GetAllWorkFlow(wfr.User)
// 		if err != nil {
// 			flogs.Errorf("Error getting workflow: %v", err)
// 			errorResponse(w, err)
// 			return
// 		}

// 		createWorkFlowResponse(w, http.StatusOK, types.WorkFlowResponse{
// 			Code:  http.StatusOK,
// 			Msg:   "success",
// 			Flows: flows,
// 		})
// 	}
// }

// func workflow() (*fworkflow.WorkflowClient, error) {
// 	return fworkflow.NewWorkflowClient(
// 		driver.WeaviateClientConf{
// 			Host:   os.Getenv(types.EnvWeaviateHost),
// 			Schema: os.Getenv(types.EnvWeaviateSchema),
// 			Key:    os.Getenv(types.EnvWewaviateKey),
// 		},
// 		driver.MongoCliConf{
// 			Uri:        os.Getenv(types.EnvMONGOHOST),
// 			DB:         os.Getenv(types.EnvMONGODB),
// 			Collection: os.Getenv(types.EnvMONGOCOLLECTION),
// 		},
// 	)
// }

// func createWorkFlowResponse(w http.ResponseWriter, code int, data types.WorkFlowResponse) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)

//		json.NewEncoder(w).Encode(data)
//	}
//
// WorkFlowHandler handles the /v1/workflow route
func WorkFlowHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize MongoDB store
	mongoStore := workflow.NewMongoStore(
		os.Getenv(types.EnvMONGOHOST),
		os.Getenv(types.EnvMONGODB),
	)

	// Initialize WorkFlowService
	service := workflow.NewWorkFlowService(mongoStore)

	// Initialize APIHandler
	apiHandler := workflow.NewAPIHandler(service)

	// Handle the API request
	apiHandler.HandleWorkFlowRequest(w, r)
}
