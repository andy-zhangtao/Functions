package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// WorkFlow represents a workflow entity
type WorkFlow struct {
	ID      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Action  string `json:"action" bson:"action"`
	StepIDs []int  `json:"step_ids" bson:"step_ids"`
}

// Step represents a step entity within a workflow
type Step struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
}

// Result represents the result of a workflow execution
type Result struct {
	WorkFlowID  string                 `json:"workflow_id"`
	Status      string                 `json:"status"`
	StepResults map[string]interface{} `json:"step_results"`
}

// WorkFlowModel represents a workflow model.
type WorkFlowModel struct {
	// WorkFlowId is the unique identifier of the workflow.
	WorkFlowId int `json:"workflow_id" bson:"workflow_id"`
	// User is the user who created the workflow.
	User string `json:"user" bson:"user"`
	// Name is the name of the workflow.
	Name string `json:"name" bson:"name"`
	// StepIds is a list of step IDs in the workflow.
	StepIds []int `json:"step_ids" bson:"step_ids"`
	// Flows is a list of flow models in the workflow.
	Flows []FlowModel `json:"flows" bson:"flows"`
}

type FlowModel struct {
	StepId   int                    `json:"step_id" bson:"step_id"`
	Name     string                 `json:"name" bson:"name"`
	Desc     string                 `json:"desc" bson:"desc"`
	Kind     string                 `json:"kind" bson:"kind"`
	Input    map[string]interface{} `json:"input" bson:"input"`
	Output   map[string]interface{} `json:"output" bson:"output"`
	ObjectID primitive.ObjectID     `json:"objectId" bson:"_id,omitempty"`
}

type WorkFlowRequest struct {
	Action   int    `json:"action"`
	User     string `json:"user"`
	Name     string `json:"name"`
	Question string `json:"question"`
}

const (
	WorkFlowActionNew = iota
	WorkFlowActionGet
	WorkFlowActionExecute
)

type WorkFlowResponse struct {
	Version string           `json:"version"`
	Msg     string           `json:"msg"`
	Code    int              `json:"code"`
	Flows   []*WorkFlowModel `json:"flows,omitempty"`
}

const (
	MongoDBWorkFlow = "workflows"
	MongoDBSteps    = "steps"
	MongoDBPlugins  = "plugins"
)

type WorkFlowBaseInfo struct {
	User string
}
