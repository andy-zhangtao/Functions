package plugins

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andy-zhangtao/Functions/tools/tgpt"
	"github.com/andy-zhangtao/Functions/tools/tplugins"
	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GPT struct {
	traceId string
	plugin  types.Plugin
	c       GPTConfig

	nextPlugin types.Plugin
	wfc        *types.WorkflowContext

	getPluginWithID func(id int) ([]types.Plugin, error)
}

type GPTConfig struct {
	Url          string  `json:"url"`
	SKey         string  `json:"skey"`
	SystemPrompt string  `json:"system_prompt"`
	Model        string  `json:"model"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float64 `json:"temperature"`
}

func NewGPTPlugin(c GPTConfig, fc *types.WorkflowContext) *GPT {
	traceId := ""
	_traceId := fc.Get(types.TraceID)
	if _traceId != nil {
		traceId = _traceId.(string)
	}

	g := &GPT{
		traceId: traceId,
		c:       c,
		wfc:     fc,
	}

	g.log("GPT plugin initialized with [%+v]", c)
	return g
}

func (p *GPT) log(format string, args ...interface{}) {
	format = "[GPT-FC-Plugin]-[info]: %s " + format
	args = append([]interface{}{p.traceId}, args...)
	logrus.Infof(format, args)
}

func (p *GPT) error(format string, args ...interface{}) {
	format = "[GPT-FC-Plugin]-[error]: %s " + format
	args = append([]interface{}{p.traceId}, args...)
	logrus.Infof(format, args)
}

func (p *GPT) Initialize(plugin types.Plugin) error {

	p.log("GPT plugin initialized with [%+v]", plugin)
	p.plugin = plugin

	input, err := p.parseGPTPlugin(plugin)
	if err != nil {
		return errors.WithMessage(err, "parse input error")
	}

	p.c.SystemPrompt = input.Prompt.System
	p.c.Model = input.Model
	p.c.MaxTokens = input.MaxTokens
	p.c.Temperature = input.Temperature

	getPluginWithID := p.wfc.Get(types.GetPluginWithID)
	if getPluginWithID == nil {
		return errors.New("get plugin with id error")
	}

	p.getPluginWithID = getPluginWithID.(func(id int) ([]types.Plugin, error))

	p.log("GPT plugin initialized with [%+v]", p.c)

	return nil
}

func (p *GPT) Execute(ctx *types.WorkflowContext, question string) error {
	p.log("GPT plugin execute with question: %s", question)

	response, err := p.do(question)
	if err != nil {
		p.error("do gpt error: %v", err)
		return errors.WithMessage(err, "do gpt error")
	}

	p.log("GPT plugin execute with response: %+v", response)

	choice := response.Choices[0]
	if strings.Contains(choice.Message.Content, "openai response error:") {
		// 如果返回的消息中包含openai response error:，则说明预测出错了
		return fmt.Errorf("%s", choice.Message.Content)
	}

	if choice.FinishReason == types.OpenAIStop &&
		choice.Message.FunctionCall == nil {
		return errors.Errorf("stop and function call is nil")
	}

	if choice.FinishReason == types.OpenAILength {
		return errors.Errorf("too length, the limit is %d, but now I has generate %d ", p.c.MaxTokens, len(choice.Message.Content))
	}

	result := make(map[string]interface{})
	// If parse success ,then fill up the result with down plugin result
	if choice.Message.FunctionCall != nil {
		pm, err := tgpt.ParseFCArgumentsToMap(choice.Message.FunctionCall.Arguments)
		if err != nil {
			return errors.WithMessage(err, "parse function call arguments error")
		}

		for k, v := range pm {
			result[k] = v
		}
	}

	// Fill up the result with the content
	p.wfc.Set(tplugins.PluginNameInChain(p.nextPlugin.Name), result)
	p.log("GPT next plugin with input: %+v", result)
	return nil
}

func (p *GPT) Finalize() (*types.WorkflowContext, error) {
	p.log("GPT plugin finalize")
	return p.wfc, nil
}

func (p *GPT) messages(question string) []types.OpenAIMessage {
	return []types.OpenAIMessage{
		{
			Role:    "system",
			Content: p.c.SystemPrompt,
		},
		{
			Role:    "user",
			Content: question,
		},
	}
}

func (p *GPT) functingCalling() ([]types.OpenAIFunction, error) {
	// 一般来说GPT模块应该是工作流第一个模块，所以这里不需要获取up plugin
	if len(p.plugin.Reference.Down) == 0 {
		return nil, nil
	}

	// 如果存在down plugin，那么就获取down plugin
	// 目前仅支持一个down plugin
	downPluginKey := p.plugin.Reference.Down[0]
	downPlugins, err := p.getPluginWithID(downPluginKey)
	if err != nil {
		return nil, errors.WithMessage(err, "getPluginWithID error")
	}

	if len(downPlugins) == 0 {
		return nil, errors.Errorf("not find plugin with %d", downPluginKey)
	}

	p.nextPlugin = downPlugins[0]
	return p.generateOpenAIFunctionViaPlugin(downPlugins[0])
}

// generateOpenAIFunctionViaPlugin 通过plugin生成OpenAIFunction
// 通过Plugin的name和describe填充OpenAIFunction数据，重点是对Plugin Input 的提取和峰值
// plugin的示例如下:
// {}
func (p *GPT) generateOpenAIFunctionViaPlugin(plugin types.Plugin) (result []types.OpenAIFunction, err error) {
	// Initialize OpenAIFunction
	of := types.OpenAIFunction{
		Name:        plugin.Name,
		Description: plugin.Descript,
	}

	// Check if Plugin Input is a map
	inputMap := plugin.Input

	// Initialize OpenAIFunctionParameters
	of.Parameters = types.OpenAIFunctionParameters{
		Type:       "object",
		Properties: make(map[string]types.OpenAiPropertyDescription),
		Required:   []string{},
	}

	// Loop through Plugin Input to populate OpenAIFunctionParameters
	for _, param := range inputMap {
		// Determine the type of the property
		valueType := param.Value.Type

		// Create OpenAiPropertyDescription
		property := types.OpenAiPropertyDescription{
			Type:        valueType,
			Description: param.Value.Description,
		}

		// Add to OpenAIFunctionParameters
		of.Parameters.Properties[param.Name] = property
		of.Parameters.Required = append(of.Parameters.Required, param.Name)
	}

	return []types.OpenAIFunction{of}, nil
}

func (p *GPT) do(question string) (res types.OpenAIResponse, err error) {
	reqModel := types.OpenAIWithFunctionRequest{
		Model:       p.c.Model,
		MaxTokens:   p.c.MaxTokens,
		Temperature: p.c.Temperature,
		Messages:    p.messages(question),
		// FunctionCall: &gi.functionName,
	}

	fc, err := p.functingCalling()
	if err != nil {
		return res, errors.WithMessage(err, "generate function calling error")
	}

	p.log("generate function calling: %+v", fc)
	reqModel.Functions = fc
	if len(fc) == 1 {
		// 如果只有一个function，那么就直接调用function
		// 如果有多个function，那么就auto
		reqModel.FunctionCallName = types.OpenAIFunctionCallName{
			Name: fc[0].Name,
		}
		p.log("only one function, invoke function call: %+v", reqModel.FunctionCallName)
	}
	if len(fc) > 1 {
		reqModel.FunctionCallName = types.OpenAIFunctionCALLAuto
		p.log("more than one function, invoke function call: %+v", reqModel.FunctionCallName)
	}

	requestBody, err := json.Marshal(reqModel)
	if err != nil {
		return res, errors.WithMessagef(err, "marshal request body error [%+v]", reqModel)
	}

	p.log("invoke gpt request: %s", string(requestBody))

	req, err := http.NewRequest(http.MethodPost, p.c.Url, bytes.NewBuffer(requestBody))
	if err != nil {
		return res, errors.WithMessagef(err, "new request error [%s]", p.c.Url)
	}

	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Send the HTTP request
	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return res, errors.WithMessagef(err, "do request error [%s]", p.c.Url)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, errors.WithMessage(err, "read response body error")
	}

	p.log("invoke gpt response: %s", string(data))

	accResponse := types.OpenAIResponse{}

	err = json.Unmarshal(data, &accResponse)
	if err != nil {
		return res, errors.WithMessagef(err, "unmarshal response body error [%s]", string(data))
	}

	return accResponse, nil
}

// parseGPTPlugin 解析输入
func (p *GPT) parseGPTPlugin(plugin types.Plugin) (input types.PluginGPTInput, err error) {

	input = types.PluginGPTInput{}
	for _, v := range plugin.Input {
		switch v.Name {
		case "prompt":
			_prompt := v.Value.Description
			if _prompt == "" {
				return input, errors.New("invalid prompt")
			}
			input.Prompt.System = _prompt
		case "max_tokens":
			_maxTokens := v.Value.Description
			if _maxTokens == "" {
				return input, errors.New("invalid max_tokens")
			}
			tokens, _ := strconv.Atoi(_maxTokens)
			input.MaxTokens = tokens
		case "temperature":
			_temperature := v.Value.Description
			if _temperature == "" {
				return input, errors.New("invalid temperature")
			}
			temperature, _ := strconv.ParseFloat(_temperature, 64)
			input.Temperature = temperature
		case "model":
			_model := v.Value.Description
			if _model == "" {
				return input, errors.New("invalid model")
			}
			input.Model = _model
		default:
			continue
		}
	}
	return input, nil
}
