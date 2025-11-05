package ticket

import (
	"time"
)

type Ticket struct {
	ID          string    `json:"id"`
	Type        Type      `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CriticalPath bool     `json:"criticalpath"`
	Status      Status    `json:"status"`
	Priority    Priority  `json:"priority"`
	CreatedBy   *string   `json:"created_by,omitempty"`
	AssignedTo  *string   `json:"assigned_to,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Files       []string  `json:"files,omitempty"`
	Comments    []string  `json:"comments,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Type string

const (
	TypeBug     Type = "bug"
	TypeFeature Type = "feature"
	TypeTask    Type = "task"
)

func (t Type) Valid() bool {
	switch t {
	case TypeBug, TypeFeature, TypeTask:
		return true
	}
	return false
}

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in-progress"
	StatusClosed     Status = "closed"
)

// Valid returns true if the status is valid
func (s Status) Valid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusClosed:
		return true
	}
	return false
}

type CriticalPath bool

const (
	CriticalPathTrue CriticalPath = true
	CriticalPathFalse CriticalPath = false
)

// Priority represents the priority level of a ticket
type Priority string

const (
	PriorityUndefined Priority = "undefined"
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Valid returns true if the priority is valid
func (p Priority) Valid() bool {
	switch p {
	case PriorityUndefined, PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}

// Link represents a connection between two tickets
type Link struct {
	FromTicketID string    `json:"from_ticket_id"`
	ToTicketID   string    `json:"to_ticket_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// Filters for querying tickets
type Filters struct {
	Status     *Status
	Type       *Type
	Priority   *Priority
	AssignedTo *string
	Tags       []string
}
