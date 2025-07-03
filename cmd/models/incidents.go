package models

import (
	"time"
)

type Incident struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    string     `json:"severity"` // Low, Medium, High
	Status      string     `json:"status"`   // Open, In Progress, Resolved
	CreatedAt   time.Time  `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}
