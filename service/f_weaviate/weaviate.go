package fweaviate

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andy-zhangtao/Functions/types"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateClient struct {
	client *weaviate.Client
}

func NewWeaviateClient(host, schema, key string) (*WeaviateClient, error) {
	cfg := weaviate.Config{
		Host:       host,
		Scheme:     schema,
		AuthConfig: auth.ApiKey{Value: key},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create weaviate client: %v", err)
	}

	return &WeaviateClient{
		client: client,
	}, nil
}

func (wc *WeaviateClient) AddNewRecord(class string, properties map[string]interface{}) (*data.ObjectWrapper, error) {
	data := make(map[string]interface{})

	for key, val := range properties {
		data[key] = val
	}

	created, err := wc.client.Data().Creator().WithClassName(class).WithProperties(data).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not create record: %v", err)
	}

	logrus.Infof("Created record with id [%+v]", created)
	return created, nil
}

// GetRecords get records
// @Summary get records
// @Description get records via filter
// Query is the filter condition
func (wc *WeaviateClient) GetRecords(class string, query types.DirayQueryModel) (results types.DirayQueryResponse, err error) {

	var operands []*filters.WhereBuilder

	userWherefilter := filters.Where()
	userWherefilter.WithPath([]string{"user"}).WithOperator(filters.Equal).WithValueText(query.User)

	operands = append(operands, userWherefilter)

	if query.Start != "" {
		startWherefilter := filters.Where()
		start, _ := strconv.ParseInt(query.Start, 10, 64)

		startWherefilter.WithPath([]string{"date"}).WithOperator(filters.GreaterThanEqual).WithValueInt(start)
		operands = append(operands, startWherefilter)
	}

	if query.End != "" {
		endWherefilter := filters.Where()
		end, _ := strconv.ParseInt(query.End, 10, 64)

		endWherefilter.WithPath([]string{"date"}).WithOperator(filters.LessThanEqual).WithValueInt(end)
		operands = append(operands, endWherefilter)
	}

	where := filters.Where().WithOperator(filters.And).WithOperands(operands)

	content := graphql.Field{Name: "content"}
	date := graphql.Field{Name: "date"}

	_additional := graphql.Field{
		Name: "_additional", Fields: []graphql.Field{
			{Name: "distance"}, // always supported
		},
	}

	ignoreDistance := true
	filterCondition := wc.client.GraphQL().Get().WithClassName(class).WithWhere(where).WithFields(content, date, _additional)
	if len(query.Keys) > 0 {
		ignoreDistance = false
		text := graphql.NearTextArgumentBuilder{}
		text.WithConcepts(query.Keys)
		filterCondition.WithNearText(&text)
	}

	response, err := filterCondition.Do(context.Background())
	if err != nil {
		return results, fmt.Errorf("could not get records: %v", err)
	}

	var contentData []string
	for _, item := range response.Data {
		// logrus.Infof("Get record [%s] with id [%+v]", name, item)
		_content, err := wc.parser(class, item, ignoreDistance)
		if err != nil {
			logrus.Errorf("could not parse object: %v", err)
			continue
		}

		contentData = append(contentData, _content...)
	}

	return types.DirayQueryResponse{
		Records: contentData,
	}, nil
}

func (wc *WeaviateClient) parser(class string, object models.JSONObject, ignoreDistance bool) (result []string, err error) {
	m, ok := object.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse object")
	}

	for key, val := range m {
		if key == class {
			// logrus.Infof("key: %s, val: %s", key, reflect.TypeOf(val))
			v, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("could not parse object")
			}

			for _, val := range v {
				m, ok := val.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("could not parse object")
				}

				if !ignoreDistance {
					content := ""
					date := ""

					for key, val := range m {
						switch key {
						case "content":
							v, ok := val.(string)
							if !ok {
								return nil, fmt.Errorf("could not parse object %+v", val)
							}
							content = v
						case "date":
							v, ok := val.(float64)
							if !ok {
								return nil, fmt.Errorf("could not parse object %+v", val)
							}
							date = time.Unix(int64(v), 0).Format("2006-01-02")
						}
					}

					if content != "" && date != "" {
						result = append(result, fmt.Sprintf("%s记录的内容是  \\%s", date, content))
					}

				} else {
					// get distance if not ignore
					if _additional, ok := m["_additional"].(map[string]interface{}); ok {
						if distance, ok := _additional["distance"].(float64); ok {
							// logrus.Infof("distance: %f", distance)
							if distance <= 0.25 {
								content := ""
								date := ""
								for key, val := range m {
									switch key {
									case "content":
										v, ok := val.(string)
										if !ok {
											return nil, fmt.Errorf("could not parse object %+v", val)
										}
										content = v
									case "date":
										v, ok := val.(float64)
										if !ok {
											return nil, fmt.Errorf("could not parse object %+v", val)
										}
										date = time.Unix(int64(v), 0).Format("2006-01-02")
									}
								}
								if content != "" && date != "" {
									result = append(result, fmt.Sprintf("%s记录的内容是  \\%s", date, content))
								}
							}
						}
					}
				}

			}
		}
	}
	return
}
