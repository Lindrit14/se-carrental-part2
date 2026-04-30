package notifier

import (
	"context"
	"log/slog"
)

// Mock writes notification intent to stdout. Used in all non-production envs.
type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Send(_ context.Context, notifType, detail string) {
	slog.Info("[NOTIFICATION]", "type", notifType, "detail", detail)
}
