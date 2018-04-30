package cmd

import (
	"io/ioutil"
	"os"

	"github.com/frozzare/max/internal/config"
	"github.com/pkg/errors"
)

func readConfig(path string) (*config.Config, error) {
	var c *config.Config
	var err error

	fi, err := os.Stdin.Stat()
	if fi.Mode()&os.ModeNamedPipe != 0 {
		buf, err := ioutil.ReadAll(os.Stdin)

		if err == nil {
			c, err = config.ReadContent(string(buf))

			if err != nil {
				return nil, errors.Wrap(err, "max")
			}
		}
	} else {
		c, err = config.ReadFile(path)
	}

	if err != nil {
		return nil, errors.Wrap(err, "max")
	}

	if c == nil {
		return nil, errors.New("max: bad config")
	}

	return c, nil
}
