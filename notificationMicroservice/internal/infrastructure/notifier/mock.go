package notifier

import (
	"context"
	"log/slog"
)

// Mock implements Notifier by writing structured logs to stdout.
// Used for local dev and tests — no real email is sent.
type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Send(_ context.Context, msg Message) error {
	slog.Info("[NOTIFICATION]",
		"to", redactEmail(msg.To),
		"subject", msg.Subject,
		"text_preview", preview(msg.TextBody, 200))
	return nil
}

func preview(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
