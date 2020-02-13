package config

import (
	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
)

// Config is the configuration.
type Config struct {
	Loglevel      string `arg:"--loglevel,env:LOGLEVEL"`
	ListenAddress string `arg:"--listen,env:LISTEN_ADDRESS"`
	KubeConfig    string `arg:"--kubeconfig,env:KUBECONFIG"`
	Namespace     string `arg:"--namespace,env:NAMESPACE"`
	IndexHTML     string `arg:"--index-html,env:INDEX_HTML"`
}

// New creates a new Config by parsing the environment and flags.
func New() (Config, error) {
	c := Config{
		ListenAddress: ":8080",
	}

	if err := arg.Parse(&c); err != nil {
		return c, errors.Wrap(err, "failed to parse config")
	}
	return c, nil
}
