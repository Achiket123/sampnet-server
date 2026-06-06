package calendar

import (
	"time"
)

type CalendarEvent struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	EventType      string    `json:"event_type"` // "task", "milestone", "leave", "attendance"
	EntityID       uint      `json:"entity_id"`
	Status         string    `json:"status"`
	AssignedToName string    `json:"assigned_to_name,omitempty"`
}