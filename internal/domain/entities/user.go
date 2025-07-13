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

// Validate ...
func (u *User) Validate() error {
	if u.ID == 0 {
		return ErrInvalidUserID
	}

	if u.UserName == "" {
		return ErrInvalidUserName
	}

	return nil
}
