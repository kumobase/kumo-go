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
)
