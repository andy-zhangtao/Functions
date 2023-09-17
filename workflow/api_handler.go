package workflow

// APIHandler 结构体，用于处理API请求。
// NewAPIHandler 函数，用于初始化 APIHandler。
// HandleWorkFlowRequest 函数，用于处理 /v1/workflow API端点。这个函数会读取工作流ID（假设它是作为查询参数传递的），执行工作流，并返回序列化的结果。

import (
	"encoding/json"
	"net/http"
)

// APIHandler handles the API requests for workflows
type APIHandler struct {
	Service *WorkFlowService
}

// NewAPIHandler initializes a new APIHandler
func NewAPIHandler(service *WorkFlowService) *APIHandler {
	return &APIHandler{Service: service}
}

// HandleWorkFlowRequest handles the /v1/workflow API endpoint
func (handler *APIHandler) HandleWorkFlowRequest(w http.ResponseWriter, r *http.Request) {
	// For demonstration, we assume the workflow ID is passed as a query parameter
	workflowID := r.URL.Query().Get("id")

	// Execute the workflow
	result, err := handler.Service.ExecuteWorkFlow(workflowID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Serialize and return the result
	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to serialize result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
