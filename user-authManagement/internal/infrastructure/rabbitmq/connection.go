// Package rabbitmq is the AMQP adapter that implements ports.EventPublisher.
package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection wraps an AMQP connection + confirm-mode channel and
// transparently reconnects on failure.
type Connection struct {
	url      string
	exchange string
	backoff  time.Duration
	logger   *slog.Logger

	mu      sync.RWMutex
	conn    *amqp.Connection
	channel *amqp.Channel
	closed  bool
	confirms chan amqp.Confirmation
}

func Dial(ctx context.Context, url, exchange string, backoff time.Duration, logger *slog.Logger) (*Connection, error) {
	c := &Connection{url: url, exchange: exchange, backoff: backoff, logger: logger}
	if err := c.connect(); err != nil {
		return nil, err
	}
	go c.watch(ctx)
	return c, nil
}

func (c *Connection) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("amqp dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("amqp channel: %w", err)
	}
	if err := ch.ExchangeDeclare(c.exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("declare exchange: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("confirm mode: %w", err)
	}
	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.confirms = ch.NotifyPublish(make(chan amqp.Confirmation, 16))
	c.mu.Unlock()
	return nil
}

// watch reconnects in a loop until the parent context is cancelled.
func (c *Connection) watch(ctx context.Context) {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()
		if conn == nil {
			return
		}
		notify := conn.NotifyClose(make(chan *amqp.Error, 1))
		select {
		case <-ctx.Done():
			c.Close()
			return
		case err, ok := <-notify:
			if !ok {
				return
			}
			if c.isClosed() {
				return
			}
			c.logger.Warn("rabbitmq connection lost", slog.Any("error", err))
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(c.backoff):
				}
				if err := c.connect(); err != nil {
					c.logger.Error("rabbitmq reconnect failed", slog.Any("error", err))
					continue
				}
				c.logger.Info("rabbitmq reconnected")
				break
			}
		}
	}
}

func (c *Connection) Channel() *amqp.Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channel
}

func (c *Connection) Exchange() string { return c.exchange }

func (c *Connection) Healthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.conn.IsClosed()
}

func (c *Connection) isClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
