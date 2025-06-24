package utils

type History struct {
	values  []string
	current int
}

// NOTE: consider ddl.

func NewHistory() *History {
	return &History{}
}

func (h *History) Save(val string) {
	h.values = append(h.values, val)
	h.current = 0
}

func (h *History) Prev() string {
	if len(h.values) == 0 {
		return ""
	}

	if h.current < len(h.values) {
		h.current++
	}

	return h.values[len(h.values)-h.current]
}

func (h *History) Next() string {
	if h.current > 0 {
		h.current--
	}

	if h.current == 0 {
		return ""
	}

	return h.values[len(h.values)-h.current]
}
