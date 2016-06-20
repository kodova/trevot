package trevot

type Server struct {
	a   Adapter
	mux *messageMux
}

func Listen(adapter Adapter) {
	defaultServer := &Server{a: adapter, mux: newMessageMux()}
	defaultServer.listen()
}

var defaultServer *Server;

func (s *Server) listen() {
	for {
		msg := <-s.a.Messages()
		go func() {
			ctx := MessageContext{Channel: msg.Channel, Sender: msg.User, Addressed: false}
			handler := s.mux.Handler(ctx)
			handler.HandleMessage(&channelResponder{a: s.a, c: msg.Channel}, ctx)
		}()
	}
}

func (s *Server) newMessageResponder(msg *Message) Responder {
	return
}

type channelResponder struct {
	a Adapter
	c *Channel
}

func (cr *channelResponder) Reply(message string) {
	cr.a.Send(cr.c, message)
}
