package types

import (
	"time"
)

type Alert struct {
	status    string
	labels    string
	timestamp uint
}
