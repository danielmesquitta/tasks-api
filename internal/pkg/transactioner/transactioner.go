package transactioner

import (
	"context"
)

var (
	CtxTxKey = struct{}{}
)

type Transactioner interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
