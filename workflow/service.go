package workflow

// WorkFlowService 结构体，用于处理工作流的主要逻辑。
// NewWorkFlowService 函数，用于初始化 WorkFlowService。
// ExecuteWorkFlow 函数，用于执行工作流。这个函数会根据工作流ID读取工作流，检查动作是否为 "execute"，读取并执行步骤，并最终返回结果。

import (
	"os"

	"github.com/andy-zhangtao/Functions/plugins"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

// WorkFlowService is the main service for handling workflows
type WorkFlowService struct {
	Store     *MongoStore
	traceId   string
	pluginMap map[string]plugins.Plugin
	ctx       *types.WorkflowContext
}

// NewWorkFlowService initializes a new WorkFlowService
func NewWorkFlowService(store *MongoStore, traceId string) *WorkFlowService {

	wfs := &WorkFlowService{Store: store, traceId: traceId}
	wfs.initContext()

	wfs.pluginMap = map[string]plugins.Plugin{
		"weaviate-function-calling": plugins.NewGPTPlugin(plugins.GPTConfig{
			Url:  "https://api.openai.com/v1/chat/completions",
			SKey: os.Getenv(types.PluginGPTSKey),
		}, wfs.ctx),
		"weaviate": plugins.NewWeaviatePlugin(plugins.WeaviateConfig{
			C: weaviate.Config{
				Host:       os.Getenv(types.EnvWeaviateHost),
				Scheme:     os.Getenv(types.EnvWeaviateSchema),
				AuthConfig: auth.ApiKey{Value: os.Getenv(types.EnvWewaviateKey)},
			},
		}, wfs.ctx),
	}
	return wfs
}

func (service *WorkFlowService) initContext() {
	service.log("initContext")
	service.ctx = types.NewWorkFlowContext()

	service.ctx.Set(types.TraceID, service.traceId)
	service.ctx.Set(types.GetPluginWithID, service.Store.GetPluginByPluginKey)
	service.log("initContext done")
}

func (service *WorkFlowService) log(format string, args ...interface{}) {
	format = "[WorkFlowService]-[info]-[%s] " + format
	args = append([]interface{}{service.traceId}, args...)
	logrus.Infof(format, args...)
}

func (service *WorkFlowService) error(format string, args ...interface{}) {
	format = "[WorkFlowService]-[error]-[%s] " + format
	args = append([]interface{}{service.traceId}, args...)
	logrus.Errorf(format, args...)
}

// ExecuteWorkFlow executes a workflow based on its ID
func (service *WorkFlowService) ExecuteWorkFlow(workflowID string, query types.WorkFlowRequest) (*types.Result, error) {
	// Read workflow by ID
	service.log("Executing workflow: %s", workflowID)
	workflow, err := service.Store.GetWorkFlowByID(workflowID)
	if err != nil {
		service.error("Error getting workflow: %v", err)
		return nil, errors.WithMessage(err, "error getting workflow")
	}

	// Check if the action is "execute"
	if workflow.Action != "execute" {
		return nil, errors.New("Invalid action")
	}

	service.log("Executing workflow: %+v", workflow)

	service.ctx.Set(types.CtxOriginQuery, types.WorkFlowBaseInfo{
		User: query.User,
	})
	// Execute steps and collect results
	stepResults := make(map[string]interface{})
	// Read steps by IDs
	for _, step := range workflow.StepIDs {
		plugins, err := service.Store.GetPluginByPluginKey(step)
		if err != nil {
			service.error("Error getting plugins: %v", err)
			return nil, errors.WithMessage(err, "error getting steps")
		}

		for _, plugin := range plugins {
			// Here you can add the actual execution logic for each step
			// For demonstration, we just log the step
			service.log("Executing plugin: %s(%s)", plugin.Name, plugin.Descript)

			if p, exist := service.pluginMap[plugin.Name]; exist {
				err := p.Initialize(plugin)
				if err != nil {
					service.error("plugin: %v initialize error: %v", plugin.Name, err)
					return nil, errors.WithMessage(err, "error getting plugin")
				}

				err = p.Execute(service.ctx, query.Question)
				if err != nil {
					service.error("plugin: %v execute error: %v", plugin.Name, err)
					return nil, errors.WithMessage(err, "error getting plugin")
				}

				_, err = p.Finalize()
				if err != nil {
					service.error("plugin: %v finalize error: %v", plugin.Name, err)
					return nil, errors.WithMessage(err, "error getting plugin")
				}

			} else {
				service.error("plugin: %v not exist", plugin.Name)
				return nil, errors.WithMessage(err, "error getting plugin")
			}
			stepResults[plugin.Name] = "Success" // Replace with actual result
		}
	}

	// Create and return the result
	result := &types.Result{
		WorkFlowID:  workflow.ID,
		Status:      "Completed",
		StepResults: stepResults,
	}

	return result, nil
}
