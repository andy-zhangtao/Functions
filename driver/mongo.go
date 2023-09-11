package driver

import (
	"context"
	"fmt"

	"github.com/andy-zhangtao/Functions/tools/flogs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoCli struct {
	cli        *mongo.Client
	db         string
	collection string
}

type MongoCliConf struct {
	Uri        string
	DB         string
	Collection string
}

func NewMongoCli(conf MongoCliConf) (*MongoCli, error) {

	flogs.Infof("NewMongoCli uri: %s", conf.Uri)

	clientOpts := options.Client().ApplyURI(conf.Uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo error: %w", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	return &MongoCli{
		cli:        client,
		db:         conf.DB,
		collection: conf.Collection,
	}, err
}

func (mc *MongoCli) SaveDataToMongo(data interface{}) error {
	collection := mc.cli.Database(mc.db).Collection(mc.collection)
	_, err := collection.InsertOne(context.Background(), data)
	return err
}

func (mc *MongoCli) FindWorkFlow(filter interface{}, object interface{}) (find bool, err error) {
	collection := mc.cli.Database(mc.db).Collection(mc.collection)
	err = collection.FindOne(context.Background(), filter).Decode(object)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (mc *MongoCli) FindAllWorkFlows(filter interface{}, objects interface{}) error {
	collection := mc.cli.Database(mc.db).Collection(mc.collection)
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}

	return cursor.All(context.Background(), objects)
}
