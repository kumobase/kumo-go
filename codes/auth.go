package codes

// Auth-module wire codes returned by /api/v1/auth/* endpoints.
//
// These cover registration, login (password + Google), email verification,
// and the password-reset flow. Branch on these constants rather than the
// human-readable Message field.
const (
	// InvalidCredentials — login failed because the email/password pair did
	// not match. Returned identically for an unknown email AND a wrong
	// password so the endpoint does not leak which emails have accounts.
	InvalidCredentials = "INVALID_CREDENTIALS"

	// EmailNotVerified — the account exists but its email has not been
	// confirmed yet. The client should prompt the user to check their inbox
	// (or offer to resend the verification email).
	EmailNotVerified = "EMAIL_NOT_VERIFIED"

	// NoPasswordSet — login was attempted with a password on an account that
	// has no password (it was created via Google sign-in). The client should
	// steer the user to Google sign-in.
	NoPasswordSet = "NO_PASSWORD_SET"

	// AccountLocked — login is temporarily blocked after too many consecutive
	// failed attempts on this account. Returned with HTTP 429; the user should
	// wait for the lockout window to elapse (or reset their password, which
	// also clears the lock) before retrying.
	AccountLocked = "ACCOUNT_LOCKED"

	// EmailAlreadyRegistered — registration was rejected because the email is
	// already in use (UNIQUE(email) violation).
	EmailAlreadyRegistered = "EMAIL_ALREADY_REGISTERED"

	// UserAlreadyVerified — the verification link was followed for an account
	// that is already verified. Idempotent no-op from the user's perspective.
	UserAlreadyVerified = "USER_ALREADY_VERIFIED"

	// VerificationTokenInvalid — the email-verification token in the link did
	// not match any pending account (mistyped, already consumed, or stale).
	VerificationTokenInvalid = "VERIFICATION_TOKEN_INVALID"

	// ResetTokenInvalid — the password-reset token did not match any issued
	// token.
	ResetTokenInvalid = "RESET_TOKEN_INVALID"

	// ResetTokenUsed — the password-reset token was already consumed by a
	// successful reset. Request a fresh reset email.
	ResetTokenUsed = "RESET_TOKEN_USED"

	// ResetTokenExpired — the password-reset token is past its expiry window.
	// Request a fresh reset email.
	ResetTokenExpired = "RESET_TOKEN_EXPIRED"

	// GoogleOAuthUnavailable — Google sign-in is not configured on this
	// deployment (no client id). Returned with HTTP 503.
	GoogleOAuthUnavailable = "GOOGLE_OAUTH_UNAVAILABLE"

	// GoogleTokenInvalid — the supplied Google ID token failed verification
	// (bad signature, wrong audience, or expired).
	GoogleTokenInvalid = "GOOGLE_TOKEN_INVALID"

	// GoogleEmailUnverified — the Google account's email is not verified by
	// Google, so it cannot be used to establish a Kumo account.
	GoogleEmailUnverified = "GOOGLE_EMAIL_UNVERIFIED"

	// GoogleEmailRequired — the Google ID token carried no email claim.
	GoogleEmailRequired = "GOOGLE_EMAIL_REQUIRED"

	// GoogleAccountLinked — the Google identity is already linked to a
	// different Kumo user, so it cannot be linked to this one.
	GoogleAccountLinked = "GOOGLE_ACCOUNT_LINKED"

	// AuthInvalidRequestBody — the request body could not be parsed.
	AuthInvalidRequestBody = "AUTH_INVALID_REQUEST_BODY"

	// AuthInternalError — an unexpected server-side failure in the auth flow.
	AuthInternalError = "AUTH_INTERNAL_ERROR"

	// RefreshTokenMissing — POST /api/v1/auth/refresh was called without a
	// refresh token (neither the refresh_token cookie nor a body field).
	// Returned with HTTP 400. The client should send the user to login.
	RefreshTokenMissing = "REFRESH_TOKEN_MISSING"

	// RefreshTokenInvalid — the supplied refresh token did not match any
	// issued token (mistyped, forged, or already pruned). HTTP 401; both
	// auth cookies are cleared. The client should send the user to login.
	RefreshTokenInvalid = "REFRESH_TOKEN_INVALID"

	// RefreshTokenExpired — the refresh token is past its expiry (or the
	// session hit its absolute maximum lifetime). HTTP 401; cookies cleared.
	// The user must log in again.
	RefreshTokenExpired = "REFRESH_TOKEN_EXPIRED"

	// RefreshTokenRevoked — the refresh token was explicitly revoked
	// (logout, logout-all, password reset, or account suspension). HTTP 401;
	// cookies cleared.
	RefreshTokenRevoked = "REFRESH_TOKEN_REVOKED"

	// RefreshTokenReused — an already-rotated refresh token was presented
	// outside the rotation grace window, which is treated as token theft. The
	// entire session family is revoked server-side (RFC-6819). HTTP 401;
	// cookies cleared. The user must log in again on every device.
	RefreshTokenReused = "REFRESH_TOKEN_REUSED"

	// RefreshAccountInactive — the account was suspended, locked, unverified,
	// or deleted between token issuance and refresh, so a new access token
	// cannot be minted. The session family is revoked. HTTP 401.
	RefreshAccountInactive = "REFRESH_ACCOUNT_INACTIVE"
)

// Admin auth wire codes returned by the admin user-management endpoints
// (modules/auth/admin_handler.go). These cover safety rails on mutating
// other users' accounts: promotions/demotions, suspension, verification,
// and deletion. Branch on these constants rather than the Message field.
//
// Note: the admin email-conflict case (ErrEmailAlreadyExists) reuses the
// existing EmailAlreadyRegistered code; the no-fields-to-update case reuses
// the cross-cutting ValidationFailed code. No new code is added for those.
const (
	// UserNotFound — the target user id did not resolve to an account (HTTP
	// 404). Returned by every admin user-management endpoint.
	UserNotFound = "USER_NOT_FOUND"

	// CannotDemoteSelf — an admin attempted to change their own admin status
	// (HTTP 409). Admins cannot self-demote.
	CannotDemoteSelf = "CANNOT_DEMOTE_SELF"

	// CannotRemoveLastAdmin — removing admin status would leave the platform
	// with no admins (HTTP 409).
	CannotRemoveLastAdmin = "CANNOT_REMOVE_LAST_ADMIN"

	// CannotDeleteLastAdmin — deleting this user would leave the platform with
	// no admins (HTTP 409).
	CannotDeleteLastAdmin = "CANNOT_DELETE_LAST_ADMIN"

	// CannotUnverifyUser — an admin attempted to un-verify a user who is
	// already verified, which is not permitted (HTTP 409).
	CannotUnverifyUser = "CANNOT_UNVERIFY_USER"

	// CannotSuspendSelf — an admin attempted to suspend their own account
	// (HTTP 409).
	CannotSuspendSelf = "CANNOT_SUSPEND_SELF"

	// UserHasActiveDeployments — the target user could not be deleted because
	// they still have active deployments that must be torn down first (HTTP
	// 409).
	UserHasActiveDeployments = "USER_HAS_ACTIVE_DEPLOYMENTS"
)
