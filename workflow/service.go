package workflow

// WorkFlowService 结构体，用于处理工作流的主要逻辑。
// NewWorkFlowService 函数，用于初始化 WorkFlowService。
// ExecuteWorkFlow 函数，用于执行工作流。这个函数会根据工作流ID读取工作流，检查动作是否为 "execute"，读取并执行步骤，并最终返回结果。

import (
	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// WorkFlowService is the main service for handling workflows
type WorkFlowService struct {
	Store   *MongoStore
	traceId string
}

// NewWorkFlowService initializes a new WorkFlowService
func NewWorkFlowService(store *MongoStore, traceId string) *WorkFlowService {
	return &WorkFlowService{Store: store, traceId: traceId}
}

func (service *WorkFlowService) log(format string, args ...interface{}) {
	format = "[WorkFlowService]-[info]-[%s] " + format
	args = append([]interface{}{service.traceId}, args...)
	logrus.Infof(format, args...)
}

func (service *WorkFlowService) error(format string, args ...interface{}) {
	format = "[WorkFlowService]-[error]-[%s] " + format
	args = append([]interface{}{service.traceId}, args...)
	logrus.Errorf(format, args...)
}

// ExecuteWorkFlow executes a workflow based on its ID
func (service *WorkFlowService) ExecuteWorkFlow(workflowID string) (*types.Result, error) {
	// Read workflow by ID
	service.log("Executing workflow: %s", workflowID)
	workflow, err := service.Store.GetWorkFlowByID(workflowID)
	if err != nil {
		service.error("Error getting workflow: %v", err)
		return nil, errors.WithMessage(err, "error getting workflow")
	}

	// Check if the action is "execute"
	if workflow.Action != "execute" {
		return nil, errors.New("Invalid action")
	}

	service.log("Executing workflow: %+v", workflow)
	// Read steps by IDs
	steps, err := service.Store.GetStepsByID(workflow.StepIDs)
	if err != nil {
		service.error("Error getting steps: %v", err)
		return nil, errors.WithMessage(err, "error getting steps")
	}

	// Execute steps and collect results
	stepResults := make(map[string]interface{})
	for _, step := range steps {
		// Here you can add the actual execution logic for each step
		// For demonstration, we just log the step
		service.log("Executing step: %s, Command: %s", step.Name, step.Command)
		stepResults[step.ID] = "Success" // Replace with actual result
	}

	// Create and return the result
	result := &types.Result{
		WorkFlowID:  workflow.ID,
		Status:      "Completed",
		StepResults: stepResults,
	}

	return result, nil
}
