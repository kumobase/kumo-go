package types

import (
	"testing"
	"time"
)

func TestTicketsRoundTrip(t *testing.T) {
	ts := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	assignee := uint(7)
	rating := 5

	roundTrip(t, "CreateTicketRequest", CreateTicketRequest{
		Subject:     "cannot ssh into vps",
		Description: "connection refused on port 22",
		Category:    "technical",
		Priority:    "high",
	})
	roundTrip(t, "AddMessageRequest", AddMessageRequest{Content: "still broken after reboot"})
	roundTrip(t, "UpdateTicketRequest", UpdateTicketRequest{
		Subject:  strptr("clearer subject"),
		Priority: strptr("critical"),
	})
	roundTrip(t, "RateTicketRequest", RateTicketRequest{Rating: 4, Comment: "fast help"})
	roundTrip(t, "MessageResponse", MessageResponse{
		ID: 3, Content: "we are looking into it", IsInternal: false,
		UserName: "Support", IsAdmin: true, CreatedAt: ts,
	})
	roundTrip(t, "TicketResponse", TicketResponse{
		ID: 1, DisplayID: "TKT-000002", Subject: "s", Description: "d",
		Category: "billing", Priority: "normal", Status: "resolved",
		AssignedTo: &assignee, AssignedName: "Agent A",
		ResolvedAt: &ts, CreatedAt: ts, UpdatedAt: ts,
		Rating: &rating, RatingComment: "great", RatedAt: &ts,
		Messages: []MessageResponse{{ID: 1, Content: "hi", CreatedAt: ts}},
	})
}

func strptr(s string) *string { return &s }
