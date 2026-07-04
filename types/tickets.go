package types

import "time"

// Tickets are currently sessionOnly (the full /api/v1/tickets/* surface is
// blocked from API keys because customer-pasted ticket content can include
// secrets/PII). These types are exposed for SDK symmetry — when ticket
// access relaxes in a future release, no breaking SDK change will be
// needed.

// CreateTicketRequest is the body for POST /api/v1/tickets.
//
// Category must be one of: "billing", "technical", "quota_increase",
// "general". Priority is optional and defaults server-side; allowed values
// are "low", "normal", "high", "critical".
type CreateTicketRequest struct {
	Subject     string `json:"subject"`     // 1..500 chars
	Description string `json:"description"` // 1..5000 chars
	Category    string `json:"category"`
	Priority    string `json:"priority,omitempty"`
}

// AddMessageRequest is the body for POST /api/v1/tickets/:id/messages.
type AddMessageRequest struct {
	Content string `json:"content"` // 1..10000 chars
}

// UpdateTicketRequest is the body for PATCH /api/v1/tickets/:id. All fields are
// optional pointers (PATCH semantics): only the non-nil fields are updated. The
// server rejects the edit unless the ticket is still "open" (before staff
// engages). Category/Priority accept the same enum values as CreateTicketRequest.
type UpdateTicketRequest struct {
	Subject     *string `json:"subject,omitempty"`     // 1..500 chars
	Description *string `json:"description,omitempty"` // 1..5000 chars
	Category    *string `json:"category,omitempty"`
	Priority    *string `json:"priority,omitempty"`
}

// RateTicketRequest is the body for POST /api/v1/tickets/:id/rating. Rating is a
// 1..5 CSAT score; Comment is optional. Allowed only once the ticket is
// resolved or closed; a rating locks once the ticket is closed.
type RateTicketRequest struct {
	Rating  int    `json:"rating"` // 1..5
	Comment string `json:"comment,omitempty"`
}

// TicketResponse is the detail shape returned by GET /api/v1/tickets/:id.
// Messages is populated on the detail endpoint and omitted from list rows.
type TicketResponse struct {
	ID           uint              `json:"id"`
	DisplayID    string            `json:"display_id"`
	Subject      string            `json:"subject"`
	Description  string            `json:"description"`
	Category     string            `json:"category"`
	Priority     string            `json:"priority"`
	Status       string            `json:"status"`
	AssignedTo   *uint             `json:"assigned_to,omitempty"`
	AssignedName string            `json:"assigned_name,omitempty"`
	ResolvedAt   *time.Time        `json:"resolved_at,omitempty"`
	ClosedAt     *time.Time        `json:"closed_at,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	// Rating (1..5), RatingComment and RatedAt carry the customer's CSAT
	// feedback once a resolved/closed ticket has been rated; nil/empty until then.
	Rating        *int       `json:"rating,omitempty"`
	RatingComment string     `json:"rating_comment,omitempty"`
	RatedAt       *time.Time `json:"rated_at,omitempty"`
	Messages      []MessageResponse `json:"messages,omitempty"`
}

// MessageResponse is one message inside a TicketResponse. IsAdmin lets the
// UI render staff replies differently. IsInternal is true for staff-only
// notes that are filtered out of customer-visible message lists.
type MessageResponse struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	IsInternal bool      `json:"is_internal"`
	UserName   string    `json:"user_name,omitempty"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
}
