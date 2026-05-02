package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

const (
	UserEventsExchange    = "user.events"
	BookingEventsExchange = "booking.events"
	DLXExchange           = "carrental.events.dlx"

	NotificationsQueue    = "notifications.queue"
	NotificationsDLXQueue = "notifications.queue.dlx"

	maxDeliveryAttempts = 3
)

// declareTopology idempotently creates all exchanges, queues and bindings.
func declareTopology(ch *amqp.Channel) error {
	// DLX exchange (fanout — dead letters go here)
	if err := ch.ExchangeDeclare(DLXExchange, "fanout", true, false, false, false, nil); err != nil {
		return err
	}

	// DLX queue
	if _, err := ch.QueueDeclare(NotificationsDLXQueue, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(NotificationsDLXQueue, "#", DLXExchange, false, nil); err != nil {
		return err
	}

	// Main notifications queue — after maxDeliveryAttempts nacks, message goes to DLX.
	// x-delivery-limit is a quorum-queue feature, so the queue must be declared as
	// type=quorum. Classic queues reject the arg with PRECONDITION_FAILED.
	args := amqp.Table{
		"x-queue-type":           "quorum",
		"x-dead-letter-exchange": DLXExchange,
		"x-delivery-limit":       int32(maxDeliveryAttempts),
	}
	if _, err := ch.QueueDeclare(NotificationsQueue, true, false, false, false, args); err != nil {
		return err
	}

	// Source exchanges (declared idempotently — user-auth and booking also declare them)
	for _, ex := range []string{UserEventsExchange, BookingEventsExchange} {
		if err := ch.ExchangeDeclare(ex, "topic", true, false, false, false, nil); err != nil {
			return err
		}
	}

	// Bind the notifications queue to both exchanges with wildcard
	if err := ch.QueueBind(NotificationsQueue, "#", UserEventsExchange, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(NotificationsQueue, "#", BookingEventsExchange, false, nil); err != nil {
		return err
	}

	return nil
}
