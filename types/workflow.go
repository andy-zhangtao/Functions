package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkFlowModel struct {
	User  string      `json:"user"`
	Name  string      `json:"name"`
	Flows []FlowModel `json:"flows"`
}

type FlowModel struct {
	Name     string             `json:"name"`
	ObjectID primitive.ObjectID `json:"objectId" bson:"_id,omitempty"`
}

type WorkFlowRequest struct {
	Action int    `json:"action"`
	User   string `json:"user"`
	Name   string `json:"name"`
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
