package fmongo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andy-zhangtao/Functions/tools/flogs"
	"github.com/andy-zhangtao/Functions/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoCli struct {
	cli        *mongo.Client
	db         string
	collection string
}

func NewMongoCli(uri, db, collection string) (*MongoCli, error) {

	flogs.Infof("NewMongoCli uri: %s", uri)

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo error: %w", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	return &MongoCli{
		cli:        client,
		db:         db,
		collection: collection,
	}, err
}

// SaveDataToMongo save data to mongo
// dcm: DirayCreateModel
// mask: map[string]interface{}{"key": "value"}
func (mc *MongoCli) SaveDataToMongo(dcm types.DirayCreateModel, mask map[string]interface{}) error {
	_bData := bson.M{
		"use":     dcm.User,
		"date":    dcm.DateSave,
		"content": dcm.Body,
	}

	if len(mask) > 0 {
		for k, v := range mask {
			_bData[k] = v
		}
	}

	collection := mc.cli.Database(mc.db).Collection(mc.collection)
	_, err := collection.InsertOne(context.TODO(), _bData)
	return err
}

// QueryData query data from mongo
// query: DirayQueryModel
func (mc *MongoCli) QueryData(query types.DirayQueryModel) (results []types.DirayQueryResponse, err error) {
	collection := mc.cli.Database(mc.db).Collection(mc.collection)

	_bData := bson.M{}

	if query.User != "" {
		_bData["use"] = query.User
	}

	if query.Start != "" && query.End == "" {
		_start, err := time.Parse("2006-01-02", query.Start)
		if err != nil {
			return nil, fmt.Errorf("parse start time error: %w", err)
		}
		_bData["date"] = bson.M{
			"$gte": _start,
		}
	}
	if query.Start == "" && query.End != "" {
		_end, err := time.Parse("2006-01-02", query.End)
		if err != nil {
			return nil, fmt.Errorf("parse end time error: %w", err)
		}
		_bData["date"] = bson.M{
			"$lte": _end,
		}
	}
	if query.Start != "" && query.End != "" {
		_start, err := time.Parse("2006-01-02", query.Start)
		if err != nil {
			return nil, fmt.Errorf("parse start time error: %w", err)
		}
		_end, err := time.Parse("2006-01-02", query.End)
		if err != nil {
			return nil, fmt.Errorf("parse end time error: %w", err)
		}
		_bData["date"] = bson.M{
			"$gte": _start,
			"$lte": _end,
		}
	}

	flogs.Infof("QueryData _bData: %+v", _bData)
	cur, err := collection.Find(context.Background(), _bData)
	if err != nil {
		return nil, fmt.Errorf("query mongo error: %w", err)
	}

	// var results []ft.DirayQueryResponse
	for cur.Next(context.Background()) {
		var episode bson.M
		err := cur.Decode(&episode)
		if err != nil {
			return nil, fmt.Errorf("query mongo error: %w", err)
		}

		results = append(results, types.DirayQueryResponse{
			Version: types.RequestVersionV1,
			Code:    http.StatusOK,
			Status:  "",
			Records: []string{
				fmt.Sprintf("%s %s", episode["date"], episode["body"]),
			},
		})
	}

	return results, nil
}
