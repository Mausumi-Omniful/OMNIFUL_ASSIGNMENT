package listener

import (
	"errors"
)

type Config struct {
	Name        string
	Initializer Initializer
}

type Configs []Config

func (lc Config) Validate() error {
	if len(lc.Name) == 0 {
		return errors.New("listener name must be present")
	}

	if lc.Initializer == nil {
		return errors.New("listener initializer must be present")
	}

	return nil
}
