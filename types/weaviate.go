package types

import (
	"time"
)

const (
	RequestVersionDefault = RequestVersionV1
	RequestVersionV1      = "v1"
)

const (
	EnvWeaviateHost   = "WEAVIATE_HOST"
	EnvWeaviateSchema = "WEAVIATE_SCHEMA"
	EnvWewaviateKey   = "WEAVIATE_KEY"
)

const (
	DiaryClassName = "Diary"
)

type DirayCreateModel struct {
	User     string      `json:"user"`
	Body     string      `json:"body"`
	Date     string      `json:"date,omitempty"`
	Tags     []string    `json:"tags"`
	Mask     interface{} `json:"mask,omitempty"`
	Version  string      `json:"version"`
	DateSave time.Time   `json:"-"` // not used in json
}

type DirayCreateResponse struct {
	Version string `json:"version"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

type DirayQueryModel struct {
	Version string   `json:"version"`
	User    string   `json:"user"`
	Start   string   `json:"start,omitempty"`
	End     string   `json:"end,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Keys    []string `json:"keys,omitempty"`
}

type DirayQueryResponse struct {
	Version string   `json:"version"`
	Status  string   `json:"status"`
	Code    int      `json:"code"`
	Records []string `json:"records"`
}

const (
	PluginTypeWeaviateCreateAction = "1"
	PluginTypeWeaviateQueryAction  = "2"
)
