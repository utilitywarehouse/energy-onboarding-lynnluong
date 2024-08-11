package models

import (
	"time"
)

type EnergyService struct {
	ServiceId  string
	UpdatedAt  time.Time
	OccurredAt time.Time
	State      string
}
