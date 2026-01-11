package handler

var reg = map[string]Handler{}

func Register(name string, h Handler) {
	reg[name] = h
}

func Get(name string) Handler {
	return reg[name]
}
