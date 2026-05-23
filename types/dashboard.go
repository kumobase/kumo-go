package types

// DashboardSummary is returned by GET /api/v1/dashboard/summary — a small
// aggregate of the user's current resource counts. New fields will be
// added here as the dashboard grows; existing fields will not be removed
// without a major SDK version bump.
type DashboardSummary struct {
	AppCount int64 `json:"app_count"`
}
