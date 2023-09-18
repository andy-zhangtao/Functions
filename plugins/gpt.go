package plugins

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

	"github.com/andy-zhangtao/Functions/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GPT struct {
	traceId string
	plugin  types.Plugin
	c       GPTConfig
	wfc     *types.WorkflowContext
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

	return &GPT{
		traceId: traceId,
		c:       c,
		wfc:     fc,
	}
}

func (p *GPT) log(format string, args ...interface{}) error {
	format = "[GPTFunctionCall]-[info]: %s " + format
	args = append([]interface{}{p.traceId}, args...)
	logrus.Infof(format, args)
	return nil
}

func (p *GPT) error(format string, args ...interface{}) error {
	format = "[GPTFunctionCall]-[error]: %s " + format
	args = append([]interface{}{p.traceId}, args...)
	logrus.Infof(format, args)
	return nil
}

func (p *GPT) Initialize(plugin types.Plugin) error {

	p.log("GPT plugin initialized with [%+v]", plugin)
	p.plugin = plugin

	input, err := p.parseInput(plugin)
	if err != nil {
		return errors.WithMessage(err, "parse input error")
	}

	p.c.SystemPrompt = input.Prompt.System
	p.c.Model = input.Model
	p.c.MaxTokens = input.MaxTokens
	p.c.Temperature = input.Temperature

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
	// TODO parse function calling result
	// If parse success ,then fill up the result with down plugin result
	// key = "plugin_N_input"
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
	// TODO get plugin via reference down plugins
	return nil, nil
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

// parseInput 解析输入
// input 为json格式, 看起来应该是:
// {"prompt":{"system":""},"max_tokens":1,"temperature":1.2,"model":""}
func (p *GPT) parseInput(plugin types.Plugin) (input types.PluginGPTInput, err error) {
	if plugin.Input.Type == "json" {
		// 解析json
		p.log("input value: %s", plugin.Input.Value)

		err = json.Unmarshal([]byte(plugin.Input.Value.(string)), &input)
		if err != nil {
			p.error("parse input error: %v", err)
			return input, errors.WithMessage(err, "parse input error")
		}
	}

	return input, errors.New("invalid input type, now only support json format")
}
