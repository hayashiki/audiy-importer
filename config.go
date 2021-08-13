package audiy_importer

import (
	"github.com/kelseyhightower/envconfig"
)

type SlackConf struct {
	GCSInputAudioBucket string `envconfig:"GCS_INPUT_AUDIO_BUCKET" required:"true"`
	SlackBotToken string `envconfig:"SLACK_BOT_TOKEN" required:"true"`
}

func NewSlackConf() (SlackConf, error) {
	env := SlackConf{}
	err := envconfig.Process("", &env)
	return env, err
}

