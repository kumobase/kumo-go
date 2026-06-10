package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values
// so an accidental rename is caught here before release.
func TestAuthCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"InvalidCredentials", InvalidCredentials, "INVALID_CREDENTIALS"},
		{"EmailNotVerified", EmailNotVerified, "EMAIL_NOT_VERIFIED"},
		{"NoPasswordSet", NoPasswordSet, "NO_PASSWORD_SET"},
		{"AccountLocked", AccountLocked, "ACCOUNT_LOCKED"},
		{"EmailAlreadyRegistered", EmailAlreadyRegistered, "EMAIL_ALREADY_REGISTERED"},
		{"UserAlreadyVerified", UserAlreadyVerified, "USER_ALREADY_VERIFIED"},
		{"VerificationTokenInvalid", VerificationTokenInvalid, "VERIFICATION_TOKEN_INVALID"},
		{"ResetTokenInvalid", ResetTokenInvalid, "RESET_TOKEN_INVALID"},
		{"ResetTokenUsed", ResetTokenUsed, "RESET_TOKEN_USED"},
		{"ResetTokenExpired", ResetTokenExpired, "RESET_TOKEN_EXPIRED"},
		{"GoogleOAuthUnavailable", GoogleOAuthUnavailable, "GOOGLE_OAUTH_UNAVAILABLE"},
		{"GoogleTokenInvalid", GoogleTokenInvalid, "GOOGLE_TOKEN_INVALID"},
		{"GoogleEmailUnverified", GoogleEmailUnverified, "GOOGLE_EMAIL_UNVERIFIED"},
		{"GoogleEmailRequired", GoogleEmailRequired, "GOOGLE_EMAIL_REQUIRED"},
		{"GoogleAccountLinked", GoogleAccountLinked, "GOOGLE_ACCOUNT_LINKED"},
		{"AuthInvalidRequestBody", AuthInvalidRequestBody, "AUTH_INVALID_REQUEST_BODY"},
		{"AuthInternalError", AuthInternalError, "AUTH_INTERNAL_ERROR"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
