package handlers

import "context"

// SendFn is the notifier signature handlers call to emit a notification.
type SendFn func(ctx context.Context, notifType, detail string)
