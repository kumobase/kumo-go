package codes

// Ticket-module wire codes returned by /api/v1/tickets/* and the admin ticket
// endpoints. Branch on these constants rather than the human-readable Message.
const (
	// TicketInvalidStatusTransition — the requested ticket status change is not
	// allowed from the ticket's current state (e.g. resolving an already-closed
	// ticket). Returned with HTTP 409 by both the user resolve flow and the
	// admin status-update flow.
	TicketInvalidStatusTransition = "TICKET_INVALID_STATUS_TRANSITION"
)
