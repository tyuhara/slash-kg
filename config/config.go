package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	// SLACK_SIGNING_SECRET is generated at slack side.
	// Please access to https://api.slack.com/apps and select your existing App or Click [Create New App]
	// Then click to Basic Information and move to App Credentials, you can see the Signing Secret at here.
	SlackSigningSecret string `envconfig:"SLACK_SIGNING_SECRET" required:"true"`

	// Please specify the Path where kubeconfig is put.
	KubeConfig string `envconfig:"os.Getenv("HOME") + "/.kube/config" required:"true"`
}

func ReadFromEnv() (*Env, error) {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}

	return &env, nil
}
