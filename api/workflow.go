package handler

import (
	"net/http"
	"os"

	traceid "github.com/andy-zhangtao/Functions/tools/trace_id"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/andy-zhangtao/Functions/workflow"
	"github.com/sirupsen/logrus"
)

// WorkFlowHandler handles the /v1/workflow route
func WorkFlowHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize MongoDB store
	traceId := traceid.ID()

	mongoStore := workflow.NewMongoStore(
		os.Getenv(types.EnvMONGOHOST),
		os.Getenv(types.EnvMONGODB),
		traceId,
	)

	if mongoStore == nil {
		w.Write([]byte("mongoStore is nil"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Initialize WorkFlowService
	service := workflow.NewWorkFlowService(mongoStore, traceId)

	logrus.Infof("WorkFlowHandler with %s", traceId)

	// Initialize APIHandler
	apiHandler := workflow.NewAPIHandler(service, traceId)

	logrus.Infof("HandleWorkFlowRequest with %s", traceId)

	// Handle the API request
	apiHandler.HandleWorkFlowRequest(w, r)
}
