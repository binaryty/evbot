package domain

import (
	"encoding/json"
	"time"
)

const (
	StepTitle       = "title"
	StepDescription = "description"
	StepDate        = "date"
	StepTime        = "time"
	StepMinutes     = "minutes"
	StepHours       = "hours"
	StepCompleted   = "completed"
	UNKNOWN         = "неизвестен"
)

type Event struct {
	ID          int64     `json:"ID,omitempty"`
	UserID      int64     `json:"userID,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}

type EventState struct {
	Step         string     `json:"step,omitempty"`
	TempEvent    Event      `json:"tempEvent"`
	TimePicker   TimePicker `json:"timePicker"`
	SelectedDate time.Time  `json:"selectedDate"`
	MessageID    int        `json:"messageID,omitempty"`
}

func (e *Event) Validate() error {
	if e.Title == "" {
		return ErrInvalidEventTitle
	}

	return nil
}

// ToJSON ...
func (s *EventState) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON ...
func (s *EventState) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}

	return nil
}
