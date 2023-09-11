package fworkflow

import (
	"github.com/andy-zhangtao/Functions/driver"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/andy-zhangtao/Functions/workflow"
	"github.com/pkg/errors"
)

type WorkflowClient struct {
	wf *workflow.WorkFlow
}

func NewWorkflowClient(wc driver.WeaviateClientConf, mc driver.MongoCliConf) (*WorkflowClient, error) {
	wf, err := workflow.NewWorkFlow(wc, mc)
	if err != nil {
		return nil, errors.WithMessage(err, "new workflow error")
	}

	return &WorkflowClient{
		wf: wf,
	}, nil
}

func (wc *WorkflowClient) NewWorkFlow(flow types.WorkFlowModel) (err error) {
	return wc.wf.NewWorkFlow(flow)
}

func (wc *WorkflowClient) GetWorkFlow(user string, name string) (flow *types.WorkFlowModel, err error) {
	return wc.wf.FindWorkFlow(name, user)
}

func (wc *WorkflowClient) GetAllWorkFlow(user string) (flows []*types.WorkFlowModel, err error) {
	return wc.wf.FindAllWorkFlows(user)
}
