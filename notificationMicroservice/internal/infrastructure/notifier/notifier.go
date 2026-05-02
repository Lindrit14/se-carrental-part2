package notifier

import "context"

// Message is the provider-agnostic email payload.
type Message struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
}

// Notifier is implemented by any concrete email backend (mock, ACS, …).
type Notifier interface {
	Send(ctx context.Context, msg Message) error
}
