package codes

// Ticket-module wire codes returned by /api/v1/tickets/* and the admin ticket
// endpoints. Branch on these constants rather than the human-readable Message.
const (
	// TicketInvalidStatusTransition — the requested ticket status change is not
	// allowed from the ticket's current state (e.g. resolving an already-closed
	// ticket). Returned with HTTP 409 by both the user resolve flow and the
	// admin status-update flow.
	TicketInvalidStatusTransition = "TICKET_INVALID_STATUS_TRANSITION"

	// TicketReopenNotAllowed — reopen was requested from a state that cannot be
	// reopened by the caller (e.g. a customer reopening a closed ticket, or
	// reopening a ticket that is neither resolved nor closed). HTTP 409.
	TicketReopenNotAllowed = "TICKET_REOPEN_NOT_ALLOWED"

	// TicketRatingNotAllowed — a CSAT rating was submitted while the ticket is in
	// a state that cannot be rated (only resolved/closed tickets are ratable).
	// HTTP 409.
	TicketRatingNotAllowed = "TICKET_RATING_NOT_ALLOWED"

	// TicketRatingLocked — the ticket is closed and already carries a rating, so
	// the rating can no longer be changed. HTTP 409.
	TicketRatingLocked = "TICKET_RATING_LOCKED"

	// TicketEditNotAllowed — a customer edit (subject/description/category/
	// priority) was attempted after staff engaged; edits are only allowed while
	// the ticket is still open. HTTP 409.
	TicketEditNotAllowed = "TICKET_EDIT_NOT_ALLOWED"

	// TicketInvalidAssignee — an admin tried to assign a ticket to a user that
	// does not exist or is not an admin. HTTP 409.
	TicketInvalidAssignee = "TICKET_INVALID_ASSIGNEE"
)
