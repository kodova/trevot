package trevot

import (
	"fmt"
	"regexp"
	"sync"
)

type messageMux struct {
	mu       sync.RWMutex
	handlers []muxEntry
}

type muxEntry struct {
	pattern   regexp.Regexp
	handler   Handler
	addressed bool
}

type MessageContext struct {
	Channel   *Channel
	Sender    *User
	Addressed bool
	Matches   []string
}

func newMessageMux() *messageMux {
	return &messageMux{handlers: make([]muxEntry, 0)}
}

func (mux *messageMux) Handler(m *MessageContext) Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	for _, entry := range mux.handlers {
		if entry.pattern.MatchString(m.Message.Text) {
			fmt.Printf("Regex %v, matched: %v", entry.pattern, m.Message.Text)
			return entry.handler
		}
	}
	return nil
}

func (mux *messageMux) Handle(pattern *regexp.Regexp, addressed bool, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == nil {
		panic("bot: nil regexp.Regexp")
	}
	if pattern == nil {
		panic("bot: nil handler")
	}

	mux.handlers = append(mux.handlers, muxEntry{
		pattern: *pattern,
		handler: handler,
	})
}

func (mux *messageMux) HandleFunc(pattern *regexp.Regexp, handler func(ResponseWriter, *MessageContext)) {
	mux.Handle(pattern, HandlerFunc(handler))
}
