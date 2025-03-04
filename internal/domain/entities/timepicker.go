package domain

import "time"

type TimePicker struct {
	SelectedTime time.Time
	TempHours    int
	TempMinutes  int
	Step         string
}
