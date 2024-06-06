package kafka

import (
	"context"
)

type Mock struct {
}

func (Mock) Push(_ context.Context, _ []byte) error {
	return nil
}

func (Mock) Health(ctx context.Context) error {
	return nil
}
