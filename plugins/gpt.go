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

func NewGPTPlugin(traceId string, c GPTConfig, fc *types.WorkflowContext) *GPT {
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

func (p *GPT) Initialize() error {
	p.log("GPT plugin initialized")
	return nil
}

func (p *GPT) Execute(ctx *types.WorkflowContext) error {
	return nil
}

func (p *GPT) Finalize() (*types.WorkflowContext, error) {
	return p.wfc, nil
}

func (p *GPT) messages() []types.OpenAIMessage {
	return []types.OpenAIMessage{
		{
			Role:    "system",
			Content: p.c.SystemPrompt,
		},
		{
			Role:    "user",
			Content: p.c.SKey,
		},
	}
}

func (p *GPT) functingCalling() ([]types.OpenAIFunction, error) {
	return nil, nil
}

func (p *GPT) do() (res types.OpenAIResponse, err error) {
	reqModel := types.OpenAIWithFunctionRequest{
		Model:       p.c.Model,
		MaxTokens:   p.c.MaxTokens,
		Temperature: p.c.Temperature,
		Messages:    p.messages(),
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
