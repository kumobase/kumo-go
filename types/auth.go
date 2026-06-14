package types

// Auth session DTOs. Most of the auth surface (register, login, password
// reset) is driven by the browser dashboard over HttpOnly cookies and is not
// modelled here. The refresh flow is the exception: a CLI/SDK client holding a
// session JWT can exchange a refresh token for a fresh access token without a
// full re-login, so those wire shapes are owned here.

// RefreshRequest is the body for POST /api/v1/auth/refresh.
//
// Browser clients leave RefreshToken empty and rely on the HttpOnly
// refresh_token cookie; CLI/SDK clients that do not use a cookie jar send the
// raw token (kumo_rt_…) here instead.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`

	// RememberMe is honoured only when the original login did not already set
	// it; it lets a client extend an existing session to the long-lived
	// "remember me" window on refresh. Optional; defaults to the value the
	// session was issued with.
	RememberMe bool `json:"remember_me,omitempty"`
}

// RefreshResponse is the Data payload of a successful POST
// /api/v1/auth/refresh. Browser clients can ignore both fields (the server
// also sets fresh auth_token + refresh_token cookies); CLI/SDK clients read
// them to update their stored credentials.
//
// Token is the new short-lived access JWT. RefreshToken is the rotated
// refresh token — the previous one is now invalid, so a client MUST persist
// this value and discard the old one.
type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
