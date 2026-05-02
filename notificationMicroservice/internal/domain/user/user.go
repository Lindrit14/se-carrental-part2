package user

import "time"

type User struct {
	ID        string
	Email     string
	UpdatedAt time.Time
}
