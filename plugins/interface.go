package plugins

import "github.com/andy-zhangtao/Functions/types"

// Plugin is the interface that all workflow plugins must implement
type Plugin interface {
	// Initialize is called once when the plugin is loaded
	Initialize() error

	// Execute performs the plugin's main action
	Execute(ctx *types.WorkflowContext) error

	// Finalize is called once after all workflow steps are completed
	Finalize() (*types.WorkflowContext, error)
}
