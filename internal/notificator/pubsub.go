package notificator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"test_task/internal/logger"
)

type PubSub struct {
	notificator RepositoryNotificator
	buffer      chan Notification
	logger      logger.Logger
}

func NewPubSub(notificator RepositoryNotificator, bufferSize int, l logger.Logger) *PubSub {
	return &PubSub{notificator: notificator, buffer: make(chan Notification, bufferSize), logger: l}
}

func (p *PubSub) Push(_ context.Context, data Notification) error {
	p.buffer <- data
	return nil
}

func (p *PubSub) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-p.buffer:
			byteData, err := json.Marshal(res)
			if err != nil {
				p.logger.Error(fmt.Errorf("marshal: %w", err))
				break
			}
			err = p.notificator.Push(ctx, byteData)
			if err != nil {
				p.logger.Error(fmt.Errorf("error push operationType=%v message=%v: %w", res.Type, res.Data, err))
			}
		}
	}
}

// Stop stops consumer, gives it time to send all notifications.
func (p *PubSub) Stop(cancelFunc context.CancelFunc, recheckTime, closeTimeout time.Duration) {
	ticker := time.NewTicker(recheckTime)
	defer ticker.Stop()
	timeoutTicker := time.NewTicker(closeTimeout)
	defer timeoutTicker.Stop()
	p.logger.Info("notification consumer stop started")
mainLoop:
	for {
		queueLength := len(p.buffer)
		if queueLength == 0 {
			break
		}
		select {
		case <-ticker.C:
			p.logger.Info(fmt.Sprintf("%v notifications should be processed", queueLength))
		case <-timeoutTicker.C:
			p.logger.Info(fmt.Sprintf("%v notifications should have been processed, but they weren't", queueLength))
			break mainLoop
		}
	}

	cancelFunc()

	p.logger.Info("notification consumer stop finished")
}
