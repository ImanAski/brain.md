package types

import (
	"encoding/json"

	"brain/core/handler"
	"brain/core/object"
)

type Note struct {
	Text string `json:"text"`
}

type NoteHandler struct{}

func init() {
	handler.Register("note", NoteHandler{})
}

func (NoteHandler) Index(o *object.Object, idx handler.Indexer) {
	var n Note
	_ = json.Unmarshal(o.Body, &n)
	idx.Add("note.text", n.Text, o.ID)
}

func (NoteHandler) IsTerminal(o *object.Object, _ []*object.Object) bool {
	return true
}

func (NoteHandler) Render(o *object.Object) string {
	var n Note
	_ = json.Unmarshal(o.Body, &n)
	return n.Text
}
