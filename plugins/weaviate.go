package plugins

import (
	"context"
	"reflect"

	"github.com/andy-zhangtao/Functions/tools/tplugins"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Weaviate struct {
	traceId string

	c WeaviateConfig

	plugin types.Plugin
	err    error

	wfc    *types.WorkflowContext
	action *WeaviateAction
	client *weaviate.Client

	getPluginWithID func(id int) ([]types.Plugin, error)
}

type WeaviateConfig struct {
	C weaviate.Config
}

func NewWeaviatePlugin(c WeaviateConfig, fc *types.WorkflowContext) *Weaviate {
	traceId := ""
	_traceId := fc.Get(types.TraceID)
	if _traceId != nil {
		traceId = _traceId.(string)
	}

	w := &Weaviate{
		traceId: traceId,
		c:       c,
		wfc:     fc,
	}

	client, err := weaviate.NewClient(c.C)
	if err != nil {
		w.error("could not create weaviate client: %v", err)
		w.err = err
		return w
	}

	w.client = client
	return w
}

func (p *Weaviate) log(format string, args ...interface{}) {
	format = "[Weaviate-Plugin]-[info]: %s " + format
	args = append([]interface{}{p.traceId}, args...)
	logrus.Infof(format, args...)
}

func (p *Weaviate) error(format string, args ...interface{}) {
	format = "[Weaviate-Plugin]-[error]: %s " + format
	args = append([]interface{}{p.traceId}, args...)
	logrus.Infof(format, args...)
}

func (p *Weaviate) Initialize(plugin types.Plugin) error {
	if p.err != nil {
		return p.err
	}

	p.log("Weaviate plugin initialized with [%+v]", plugin)
	p.plugin = plugin
	// get input from workflow context
	action, err := p.parseWeaviatePlugin(plugin)
	if err != nil {
		return errors.WithMessage(err, "parse weaviate plugin action error")
	}

	p.log("action %+v", action)
	return nil
}

func (p *Weaviate) Execute(ctx *types.WorkflowContext, question string) error {
	if p.err != nil {
		return p.err
	}

	data := make(map[string]interface{})

	switch p.action.action {
	case types.PluginTypeWeaviateCreateAction:
		data["title"] = p.action.data.(WeaviateModelDiary).Title
		data["body"] = p.action.data.(WeaviateModelDiary).Body
		data["tags"] = p.action.data.(WeaviateModelDiary).Tags
		data["user"] = p.action.data.(WeaviateModelDiary).User
		data["date"] = p.action.data.(WeaviateModelDiary).Date
	default:
		return errors.Errorf("action [%s] not support", p.action.action)
	}

	created, err := p.client.Data().Creator().WithClassName(types.DiaryClassName).WithProperties(data).Do(context.Background())
	if err != nil {
		return errors.WithMessage(err, "could not create record")
	}

	p.log("Created record with id [%+v]", created)
	return nil
}

func (p *Weaviate) Finalize() (*types.WorkflowContext, error) {
	if p.err != nil {
		return nil, p.err
	}

	p.log("Weaviate plugin finalized")
	return nil, nil
}

func (p *Weaviate) parseWeaviatePlugin(plugin types.Plugin) (*WeaviateAction, error) {

	inputParams := p.wfc.Get(tplugins.PluginNameInChain(plugin.Name))
	if inputParams == nil {
		return nil, errors.Errorf("plugin %s not found in workflow context", plugin.Name)
	}

	input, ok := inputParams.(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("plugin %s params not conver to map[string]interface{}, it`s a [%s] type", reflect.TypeOf(inputParams))
	}

	err := p.check(input)
	if err != nil {
		return nil, errors.WithMessage(err, "check input error")
	}

	action := p.convert(input)

	return &action, nil
}

func (p *Weaviate) check(input map[string]interface{}) error {
	if _, ok := input["action"]; !ok {
		return errors.Errorf("action not found in input")
	}

	switch input["action"] {
	case types.PluginTypeWeaviateCreateAction:
		return p.checkCreateInput(input)
	default:
		return errors.Errorf("action [%s] not support", input["action"])
	}

}

func (p *Weaviate) checkCreateInput(input map[string]interface{}) error {
	if _, ok := input["title"]; !ok {
		return errors.Errorf("title not found in input with create action")
	}

	if _, ok := input["body"]; !ok {
		return errors.Errorf("body not found in input with create action")
	}

	if _, ok := input["tags"]; !ok {
		return errors.Errorf("tags not found in input with create action")
	}

	if _, ok := input["user"]; !ok {
		return errors.Errorf("user not found in input with create action")
	}

	if _, ok := input["date"]; !ok {
		return errors.Errorf("date not found in input with create action")
	}

	return nil
}

func (p *Weaviate) convert(input map[string]interface{}) WeaviateAction {
	switch input["action"] {
	case types.PluginTypeWeaviateCreateAction:
		return p.convertCreateAction(input)
	default:
		return WeaviateAction{}
	}
}

func (p *Weaviate) convertCreateAction(input map[string]interface{}) WeaviateAction {
	return WeaviateAction{
		action: types.PluginTypeWeaviateCreateAction,
		class:  types.DiaryClassName,
		data: WeaviateModelDiary{
			Title: input["title"].(string),
			Body:  input["body"].(string),
			Tags:  input["tags"].([]string),
			User:  input["user"].(string),
			Date:  input["date"].(string),
		},
	}
}
