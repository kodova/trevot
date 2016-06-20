package trevot

import "regexp"

type Handler interface {
	HandleMessage(Responder, *MessageContext)
}

type HandlerFunc func(Responder, *MessageContext)

func (f HandlerFunc) HandleMessage(r Responder, m *MessageContext) {
	f(r, m)
}

func Eavesdrop(pattern string, handler Handler)  {
	defaultServer.mux.Handle(regexp.MustCompile(pattern), false, handler)
}

func EavesdropFunc(pattern string, handler func(Responder, *MessageContext))  {
	defaultServer.mux.Handle(regexp.MustCompile(pattern), false, HandlerFunc(handler))
}

func Addressed(pattern string, handler Handler)  {
	defaultServer.mux.Handle(regexp.MustCompile(pattern), true, handler)
}

func AddressedFunc(pattern string, handler func(Responder, *MessageContext))  {
	defaultServer.mux.Handle(regexp.MustCompile(pattern), true, HandlerFunc(handler))
}


