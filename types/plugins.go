package types

// WorkflowContext represents the shared context of a workflow

type WorkflowContext struct {
	Data map[string]interface{}
}

func NewWorkFlowContext() *WorkflowContext {
	return &WorkflowContext{
		Data: make(map[string]interface{}),
	}
}

func (ctx *WorkflowContext) Set(key string, value interface{}) {
	ctx.Data[key] = value
}

func (ctx *WorkflowContext) Get(key string) interface{} {
	return ctx.Data[key]
}

// Plugin defines the structure for a Plugin in the system
type Plugin struct {
	PluginKey  int             `bson:"plugin_key" json:"plugin_key"`
	Name       string          `bson:"name" json:"name"`
	Descript   string          `bson:"descript" json:"descript"`
	Module     string          `bson:"module" json:"module"`
	Input      PluginIO        `bson:"input" json:"input"`
	Reference  PluginReference `bson:"reference" json:"reference"`
	InvokeType string          `bson:"invoke_type" json:"invoke_type"`
	InvokeURL  string          `bson:"invoke_url" json:"invoke_url"`
	Stage      PluginIO        `bson:"stage" json:"stage"`
}

// PluginIO defines the Input/Output structure for a Plugin
type PluginIO struct {
	Name  string     `bson:"name" json:"name"`
	Value PluginType `bson:"value" json:"value"`
}

type PluginType struct {
	Type        string `bson:"type" json:"type"`
	Description string `bson:"description" json:"description"`
}

// PluginReference defines the reference structure for a Plugin
type PluginReference struct {
	Up   int   `bson:"up" json:"up"`
	Down []int `bson:"down" json:"down"`
}

const (
	CtxPluginGPT = "x-ctx-gpt-instance"
)
