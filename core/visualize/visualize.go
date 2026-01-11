package visualize

import (
	"brain/adapters/store/sqlite"
	"brain/core/object"
)

type Visualizer interface {
	Visualize(objects []*object.Object, links []sqlite.Link) error
}

var registry = make(map[string]Visualizer)

func Register(name string, v Visualizer) {
	registry[name] = v
}

func Get(name string) Visualizer {
	return registry[name]
}
