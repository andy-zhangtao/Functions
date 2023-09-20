package tplugins

// PluginNameInChain returns the plugin name in the chain
// 通过plugin的name获取plugin在chain中的参数名称
// 例如: plugin name = "doc"，那么返回的结果就是"plugin_doc_input"
// 所有的plugin通过这个方法获取在ctx中唯一的参数Key名称
func PluginNameInChain(name string) string {
	return "plugin_" + name + "_input"
}
