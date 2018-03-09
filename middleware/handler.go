package middleware


type Handler struct {
	middleware	map[string][]interface{}
}

func NewHandler() *Handler {
	return &Handler {
		middleware: make(map[string][]interface{}),
	}
}

func (h *Handler) Use(k string, v interface{})  {

	if _, ok := h.middleware[k]; !ok {
		h.middleware[k] = []interface{}{ v }
		return
	}

	h.middleware[k] = append(h.middleware[k], v)
}

func (h *Handler) Exec(k string, call func(v interface{}, next func())) {

	if list, ok := h.middleware[k]; !ok || len(list) < 1 {
		return
	}

	var (
		i = 0
		next func()
		total = len(h.middleware[k])
	)

	next = func() {
		if i >= total {
			return
		}

		call(h.middleware[k][i], next)
		i++
	}
	next()
}