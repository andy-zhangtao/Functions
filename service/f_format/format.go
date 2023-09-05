package fformat

type FormatModel struct {
	Action  string `json:"action"`
	Version string `json:"version"`
	Format  string `json:"format"`
	User    string `json:"user"`
	Tags    string `json:"tags"`
}
