package workflow

import (
	"github.com/andy-zhangtao/Functions/driver"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type WorkFlow struct {
	weaviate *driver.WeaviateClient
	mongoCli *driver.MongoCli
	traceId  string
}

func NewWorkFlow(wc driver.WeaviateClientConf, mc driver.MongoCliConf, traceId string) (*WorkFlow, error) {
	weaviateClient, err := driver.NewWeaviateClient(wc)
	if err != nil {
		return nil, errors.WithMessage(err, "new weaviate client error")
	}

	mongoCli, err := driver.NewMongoCli(mc)
	if err != nil {
		return nil, errors.WithMessage(err, "new mongo client error")
	}

	logrus.SetLevel(logrus.InfoLevel)

	return &WorkFlow{
		weaviate: weaviateClient,
		mongoCli: mongoCli,
		traceId:  traceId,
	}, nil
}

func (wf *WorkFlow) error(format string, args ...interface{}) {
	format = "[WorkFlow]-[error]-[%s] " + format
	args = append([]interface{}{wf.traceId}, args...)
	logrus.Errorf(format, args...)
}

func (wf *WorkFlow) log(format string, args ...interface{}) {
	format = "[WorkFlow]-[info]-[%s] " + format
	args = append([]interface{}{wf.traceId}, args...)
	logrus.Infof(format, args...)
}

func (wf *WorkFlow) NewWorkFlow(flow types.WorkFlowModel) (err error) {
	// 1. save to mongo
	return wf.mongoCli.SaveDataToMongo(flow)
}

func (wf *WorkFlow) FindWorkFlow(name, user string) (*types.WorkFlowModel, error) {
	// 1. find from mongo
	flow := &types.WorkFlowModel{}
	filter := bson.M{"name": name, "user": user}
	found, err := wf.mongoCli.FindWorkFlow(filter, flow)
	if err != nil {
		return nil, err
	}

	if found {
		return flow, nil
	}

	return nil, nil
}

func (wf *WorkFlow) FindAllWorkFlows(user string) ([]*types.WorkFlowModel, error) {
	var flows []*types.WorkFlowModel
	filter := bson.M{"user": user}
	err := wf.mongoCli.FindAllWorkFlows(filter, flows)
	return flows, err
}
