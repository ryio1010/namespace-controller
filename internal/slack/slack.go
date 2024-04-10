package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Client struct {
	*slack.Client
}

// NewClient creates a new Slack client.
func NewClient(token string) *Client {
	return &Client{
		Client: slack.New(token),
	}
}

func (c *Client) PostMessage(message, channel string) error {
	log.Log.Info("SlackMessage", "message", fmt.Sprintf("Posting message to Slack channel: %s, message: %s", channel, message))
	_, _, err := c.Client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return fmt.Errorf("failed to post message to Slack channel: %w", err)
	}

	return nil
}
