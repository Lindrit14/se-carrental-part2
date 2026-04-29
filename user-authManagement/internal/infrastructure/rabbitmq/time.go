package rabbitmq

import "time"

// timeFromRFC3339 parses an RFC3339 timestamp; on parse error it returns
// the zero value, which AMQP serializes as 0.
func timeFromRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
