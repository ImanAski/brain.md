package handler

import "brain/core/object"

type Indexer interface {
	Add(field, value string, id object.ID)
}

type Handler interface {
	Index(o *object.Object, idx Indexer)
	IsTerminal(o *object.Object, children []*object.Object) bool
	Render(o *object.Object) string
}
