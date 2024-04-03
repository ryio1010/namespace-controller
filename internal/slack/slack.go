package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Client struct {
	*slack.Client
	Channel string
}

// NewClient creates a new Slack client.
func NewClient(token, channel string) *Client {
	return &Client{
		Client:  slack.New(token),
		Channel: channel,
	}
}

func (c *Client) PostMessage(message string) error {
	log.Log.Info("SlackMessage", "message", fmt.Sprintf("Posting message to Slack channel: %s, message: %s", c.Channel, message))
	_, _, err := c.Client.PostMessage(c.Channel, slack.MsgOptionText(message, false))
	if err != nil {
		return fmt.Errorf("failed to post message to Slack channel: %w", err)
	}

	return nil
}
