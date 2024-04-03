package env

type Env struct {
	SlackToken string `env:"SLACK_ACCESS_TOKEN,required"`
	Channel    string `env:"SLACK_CHANNEL,required"`
}
