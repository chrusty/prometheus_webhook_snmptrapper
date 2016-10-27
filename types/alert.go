package types

import (
	"time"
)

type Alert struct {
	Address      string
	Status       string
	Annotations  map[string]string
	Labels       map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string
}
