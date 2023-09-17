package workflow

// MongoStore 结构体，用于存储MongoDB客户端。
// NewMongoStore 函数，用于初始化MongoDB存储。
// GetWorkFlowByID 函数，用于根据ID从MongoDB中获取一个WorkFlow。
// GetStepsByID 函数，用于根据一组ID从MongoDB中获取多个Step。

import (
	"context"
	"log"
	"time"

	"github.com/andy-zhangtao/Functions/tools/flogs"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// // FindWorkFlowWithuserAndId
// // @Description: find workflow with user and id
// func (wf *WorkFlow) FindWorkFlowsWithUserAndId(user string, wfId int) ([]*types.WorkFlowModel, error) {
// 	wf.log("find workflow with user: %s, id: %d", user, wfId)
// 	// 1. find from mongo
// 	flows := []*types.WorkFlowModel{}
// 	filter := bson.M{"user": user, "workflow_id": wfId}

// 	found, err := wf.mongoCli.FindWorkFlow(filter, &flows)
// 	if err != nil {
// 		wf.error("find workflow with user: %s, id: %d error: %v", user, wfId, err)
// 		return nil, errors.WithMessagef(err, "find workflow with user: %s, id: %d error", user, wfId)
// 	}

// 	if found {
// 		wf.log("find workflow with user: %s, id: %d success", user, wfId)
// 		return flows, nil
// 	}

//		wf.log("find workflow with user: %s, id: %d not found", user, wfId)
//		return nil, nil
//	}
//
// MongoStore is a struct that holds the MongoDB client
type MongoStore struct {
	Client  *mongo.Client
	db      string
	traceId string
}

func (store *MongoStore) log(format string, args ...interface{}) {
	format = "[MongoStore]-[info]-[%s] " + format
	args = append([]interface{}{store.traceId}, args...)
	logrus.Infof(format, args...)
}

func (store *MongoStore) error(format string, args ...interface{}) {
	format = "[MongoStore]-[info]-[%s] " + format
	args = append([]interface{}{store.traceId}, args...)
	logrus.Errorf(format, args...)
}

// NewMongoStore initializes a new MongoDB store
func NewMongoStore(uri, db, traceId string) *MongoStore {
	// client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	flogs.Infof("NewMongoCli uri: %s with traceId %s", uri, traceId)

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
		// return nil, fmt.Errorf("connect to mongo error: %w", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &MongoStore{
		Client:  client,
		db:      db,
		traceId: traceId,
	}
}

// GetWorkFlowByID fetches a WorkFlow by its ID from MongoDB
func (store *MongoStore) GetWorkFlowByID(id string) (*types.WorkFlow, error) {
	store.log("get workflow with id: %s", id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := store.Client.Database(store.db).Collection(types.MongoDBWorkFlow)
	var workflow types.WorkFlow

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&workflow)
	if err != nil {
		store.error("get workflow with id: %s error: %v", id, err)
		return nil, errors.WithMessage(err, "get workflow error")
	}

	return &workflow, nil
}

// GetStepsByID fetches Steps by their IDs from MongoDB
func (store *MongoStore) GetStepsByID(ids []string) ([]types.Step, error) {
	store.log("get steps with ids: %v", ids)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := store.Client.Database(store.db).Collection(types.MongoDBSteps)
	var steps []types.Step

	filter := bson.M{"id": bson.M{"$in": ids}}
	store.log("filter: %v", filter)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		store.error("get steps with ids: %v error: %v", ids, err)
		return nil, errors.WithMessage(err, "get steps error")
	}

	for cursor.Next(ctx) {
		var step types.Step
		if err := cursor.Decode(&step); err != nil {
			store.error("get steps with ids: %v error: %v", ids, err)
			return nil, errors.WithMessage(err, "decode steps error")
		}
		steps = append(steps, step)
	}

	return steps, nil
}
