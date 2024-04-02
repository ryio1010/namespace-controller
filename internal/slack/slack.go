package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	slackToken = "token"
)

type Client struct {
	*slack.Client
}

// NewClient creates a new Slack client.
func NewClient() *Client {
	return &Client{
		Client: slack.New(slackToken),
	}
}

// MessageInfo represents a message to be sent to a Slack channel.
type MessageInfo struct {
	Channel string
	Message string
}

func (c *Client) PostMessage(mi *MessageInfo) error {
	log.Log.Info("SlackMessage", "message", fmt.Sprintf("Posting message to Slack channel: %s, message: %s", mi.Channel, mi.Message))

	return nil
}
