package trevot

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
)

type slackAdapter struct {
	client *slack.Client
	rtm    *slack.RTM
	self   *slack.UserDetails
	logger *log.Logger
}

type slackResponder struct {
	a *slackAdapter
	c *Channel
}

var defaultSlackAdapter *slackAdapter

func NewSlackAdapter(token string) *Adapter {
	client := slack.New(token)
	rtm := client.NewRTM()

	defaultSlackAdapter = &slackAdapter{client: client, rtm: rtm}
	return defaultSlackAdapter
}

func (a *slackAdapter) Messages() <-chan Message {
	out := make(chan Message)
	go a.rtm.ManageConnection()
	go a.listen(out)
	return out
}

func (a *slackAdapter) Channel(id string) (*Channel, error) {
	channel, err := a.client.GetChannelInfo(id)
	if err != nil {
		return nil, err
	}

	return &Channel{
		Id:   channel.ID,
		Nmae: channel.Name,
	}
}

func (a *slackAdapter) User(id string) (*User, error) {
	user, err := a.client.GetUserInfo(id)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:   user.ID,
		Name: user.Name,
	}
}

func (a *slackAdapter) Send(c *Channel, message string) error {
	_, _, err := a.client.PostMessage(c.Id, message)
	return err
}

func (a *slackAdapter) listen(msgOut chan Message) {
	for {
		select {
		case msg := <-a.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				a.handleConnectedEvent(ev)
			case *slack.MessageEvent:
				go a.handleMessageEvent(ev, msgOut)
			case *slack.RTMError:
				panic(ev.Error())
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				panic("Failed to auth")
			}
		}
	}

}

func (a *slackAdapter) handleConnectedEvent(ev *slack.ConnectedEvent) {
	info := ev.Info
	a.self = info.User
}

func (a *slackAdapter) handleMessageEvent(ev *slack.MessageEvent, msgOut chan Message) {
	channel, err := a.Channel(ev.Channel)
	if err != nil {
		a.logger.Panicf("trevot: Unable handleMessageEvent, failed to find channel id: %v", ev.Channel)
	}

	user, err := a.User(ev.User)
	if err != nil {
		a.logger.Panicf("trevot: Unable handleMessageEvent, failed to find user id: %v", ev.User)
	}

	msgOut <- &Message{
		Text:    ev.Text,
		Channel: channel,
		User:    user,
	}
}

