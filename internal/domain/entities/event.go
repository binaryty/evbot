package domain

import "time"

const (
	StepTitle       = "title"
	StepDescription = "description"
	StepDate        = "date"
	StepTime        = "time"
	StepCompleted   = "completed"
)

type Event struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

type EventState struct {
	Step      string
	TempEvent Event
	CreatedAt time.Time
}
