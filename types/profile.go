package types

// Profile DTOs cover the read surface (GET /profile, GET /profile/balance,
// GET /profile/has-password). The mutation endpoints (PATCH /profile, password
// management) are sessionOnly today, so API-key clients cannot reach them —
// their request shapes are not exposed in this SDK.

// GetProfileResponse is returned by GET /api/v1/profile.
//
// IsAdmin is included only so consumers can render different UI; it does NOT
// expand what an API key can call (admin routes still 403 with
// API_KEY_ADMIN_FORBIDDEN regardless of this flag).
type GetProfileResponse struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	IsVerified  bool   `json:"is_verified"`
	IsAdmin     bool   `json:"is_admin"`
	HasPassword bool   `json:"has_password"`
	HasGoogle   bool   `json:"has_google"`
}

// GetBalanceResponse is returned by GET /api/v1/profile/balance. Balance is
// a decimal string in IDR (the platform's accounting currency).
type GetBalanceResponse struct {
	Balance string `json:"balance"`
}

// HasPasswordResponse is returned by GET /api/v1/profile/has-password and
// lets a dashboard decide whether to render "set password" vs "change
// password" UI without leaking other profile fields.
type HasPasswordResponse struct {
	HasPassword bool `json:"has_password"`
}
