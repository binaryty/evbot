package domain

import "time"

type User struct {
	ID        int64
	FirstName string
	UserName  string
}

type Participant struct {
	User
	RegisteredAt time.Time
}
