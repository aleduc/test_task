package notificator

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Manager struct {
	notificators []RepositoryNotificator
	eg           *errgroup.Group
}

func NewManager(notificators []RepositoryNotificator) *Manager {
	return &Manager{notificators: notificators, eg: &errgroup.Group{}}
}

func (m *Manager) Push(ctx context.Context, data Notification) error {
	for _, v := range m.notificators {
		m.eg.Go(func() error {
			res, err := json.Marshal(data)
			if err != nil {
				return err
			}
			return v.Push(ctx, res)
		})
	}
	err := m.eg.Wait()
	if err != nil {
		err = fmt.Errorf("notificator task: %w", err)
	}
	return err
}
