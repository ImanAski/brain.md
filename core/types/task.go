package types

import (
	"encoding/json"

	"brain/core/handler"
	"brain/core/object"
)

type Task struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	Due    string `json:"due,omitempty"`
}

type TaskHandler struct{}

func init() {
	handler.Register("task", TaskHandler{})
}

func (TaskHandler) Index(o *object.Object, idx handler.Indexer) {
	var t Task
	_ = json.Unmarshal(o.Body, &t)
	idx.Add("task.status", t.Status, o.ID)
}

func (TaskHandler) IsTerminal(o *object.Object, children []*object.Object) bool {
	return len(children) == 0
}

func (TaskHandler) Render(o *object.Object) string {
	var t Task
	if err := json.Unmarshal(o.Body, &t); err != nil {
		return "err"
	}
	return "[" + t.Status + "] " + t.Title
}
