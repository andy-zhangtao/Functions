package workflow

// MongoStore 结构体，用于存储MongoDB客户端。
// NewMongoStore 函数，用于初始化MongoDB存储。
// GetWorkFlowByID 函数，用于根据ID从MongoDB中获取一个WorkFlow。
// GetStepsByID 函数，用于根据一组ID从MongoDB中获取多个Step。

import (
	"context"
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

	flogs.Infof("NewMongoCli uri: %s with traceId %s", uri, traceId)

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		logrus.Errorf("connect to mongo error: %v", err)
		return nil
		// return nil, fmt.Errorf("connect to mongo error: %w", err)
	}

	logrus.Infof("Ping MongoDB!")
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logrus.Errorf("ping mongo error: %v", err)
		return nil
	}

	logrus.Infof("Connected to MongoDB!")
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

func (store *MongoStore) GetPluginByPluginKey(id int) ([]types.Plugin, error) {
	store.log("get plugin with id: %d", id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := store.Client.Database(store.db).Collection("plugins") // Assume the collection name is "plugins"
	var plugins []types.Plugin

	filter := bson.M{"plugin_key": id}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		// Logging (replace with your logging logic)
		store.error("get plugin with PluginKey: %d error: %v", id, err)
		return nil, errors.WithMessage(err, "get plugins error")
	}

	for cursor.Next(ctx) {
		var plugin types.Plugin
		if err := cursor.Decode(&plugin); err != nil {
			// Logging (replace with your logging logic)
			store.error("get plugin with PluginKey: %d error: %v", id, err)
			return nil, errors.WithMessage(err, "decode plugins error")
		}
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}
