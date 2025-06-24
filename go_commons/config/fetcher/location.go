package fetcher

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/config/model"
	"os"
)

type native struct {
	location string
}

func NewNativeFetcher(ctx context.Context, location string) (Fetcher, error) {
	nf := &native{
		location: location,
	}

	return nf, nil
}

func (nf *native) GetConfig(ctx context.Context) (*model.Config, error) {
	data, err := os.ReadFile(nf.location)
	if err != nil {
		return nil, fmt.Errorf("error while reading local file: %w", err)
	}

	configMap, err := ParseYAMLToConfigMap(string(data))
	if err != nil {
		return nil, err
	}

	return model.NewConfig(configMap), nil
}
