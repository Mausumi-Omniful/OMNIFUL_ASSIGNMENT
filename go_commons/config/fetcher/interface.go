package fetcher

import (
	"context"
	"github.com/omniful/go_commons/config/model"
)

type Fetcher interface {
	GetConfig(ctx context.Context) (*model.Config, error)
}
