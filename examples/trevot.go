package main

import "github.com/kodova/trevot"

func main() {
	trevot.Eavesdrop("", func(trevot.Responder, ) {

	})

	trevot.Addressed("", func() {

	})

	trevot.Listen(&trevot.NewSlackAdapter(""))
}

//func handleHello(r bot.ResponseWriter, m *bot.MessageContext) {
//	r.Reply("Hello, I Hear you")
//}

