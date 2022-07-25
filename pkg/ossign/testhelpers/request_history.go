package testhelpers

import (
	"net/http"
)

type RequestHistoryEntry struct {
	Request *http.Request
	Body    []byte
}

type HttpRequestHistory struct {
	Stack []*RequestHistoryEntry
}

func NewHttpRequestHistory() *HttpRequestHistory {
	return &HttpRequestHistory{
		Stack: []*RequestHistoryEntry{},
	}
}

func (h *HttpRequestHistory) Push(e *RequestHistoryEntry) {
	h.Stack = append(h.Stack, e)
}

func (h *HttpRequestHistory) Pop() *RequestHistoryEntry {
	if len(h.Stack) < 1 {
		return nil
	}

	i := len(h.Stack) - 1
	e := h.Stack[i]
	h.Stack = h.Stack[:i]

	return e
}

func (h *HttpRequestHistory) Clear() {
	h.Stack = nil
}

func (h *HttpRequestHistory) Oldest() *RequestHistoryEntry {
	if len(h.Stack) < 1 {
		return nil
	}
	return h.Stack[0]
}

func (h *HttpRequestHistory) Latest() *RequestHistoryEntry {
	if len(h.Stack) < 1 {
		return nil
	}
	return h.Stack[len(h.Stack)-1]
}
