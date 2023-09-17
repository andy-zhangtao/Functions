package workflow

// WorkFlowService 结构体，用于处理工作流的主要逻辑。
// NewWorkFlowService 函数，用于初始化 WorkFlowService。
// ExecuteWorkFlow 函数，用于执行工作流。这个函数会根据工作流ID读取工作流，检查动作是否为 "execute"，读取并执行步骤，并最终返回结果。

import (
	"log"

	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
)

// WorkFlowService is the main service for handling workflows
type WorkFlowService struct {
	Store *MongoStore
}

// NewWorkFlowService initializes a new WorkFlowService
func NewWorkFlowService(store *MongoStore) *WorkFlowService {
	return &WorkFlowService{Store: store}
}

// ExecuteWorkFlow executes a workflow based on its ID
func (service *WorkFlowService) ExecuteWorkFlow(workflowID string) (*types.Result, error) {
	// Read workflow by ID
	workflow, err := service.Store.GetWorkFlowByID(workflowID)
	if err != nil {
		return nil, err
	}

	// Check if the action is "execute"
	if workflow.Action != "execute" {
		return nil, errors.New("Invalid action")
	}

	// Read steps by IDs
	steps, err := service.Store.GetStepsByID(workflow.StepIDs)
	if err != nil {
		return nil, err
	}

	// Execute steps and collect results
	stepResults := make(map[string]interface{})
	for _, step := range steps {
		// Here you can add the actual execution logic for each step
		// For demonstration, we just log the step
		log.Printf("Executing step: %s, Command: %s", step.Name, step.Command)
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
