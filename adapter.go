package trevot

type Adapter interface {
	Messages() <-chan Message
	Channel(id string) (*Channel, error)
	User(id string) (*User, error)
	Send(c *Channel, message string) (error)
}

type Message struct {
	Text    string
	User    *User
	Channel *Channel
}

type Responder interface {
	Reply(message string)
}

type User struct {
	Id   string
	Name string
}

type Channel struct {
	Id   string
	Nmae string
}
