package domain

import "time"

type Calendar struct {
	CurrentDate time.Time
}

func NewCalendar() *Calendar {
	return &Calendar{
		CurrentDate: time.Now(),
	}
}

// NextMonth ...
func (c *Calendar) NextMonth() {
	c.CurrentDate = c.CurrentDate.AddDate(0, 1, 0)
}

// PrevMonth ...
func (c *Calendar) PrevMonth() {
	c.CurrentDate = c.CurrentDate.AddDate(0, -1, 0)
}
