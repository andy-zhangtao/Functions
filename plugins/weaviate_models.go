package plugins

type WeaviateModel struct {
	Action  string `json:"action"`
	Content string `json:"content"`
}

type WeaviateModelDiary struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
	User  string   `json:"user"`
	Date  string   `json:"date"`
}

type WeaviateAction struct {
	action string
	class  string
	data   interface{}
}
