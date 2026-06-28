package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values so an
// accidental rename is caught here before release.
func TestRunnersCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"RunnerSpecNotFound", RunnerSpecNotFound, "RUNNER_SPEC_NOT_FOUND"},
		{"RunnerJobNotFound", RunnerJobNotFound, "RUNNER_JOB_NOT_FOUND"},
		{"RunnerValidationFailed", RunnerValidationFailed, "RUNNER_VALIDATION_FAILED"},
		{"RunnerUnauthorized", RunnerUnauthorized, "RUNNER_UNAUTHORIZED"},
		{"RunnerInvalidID", RunnerInvalidID, "RUNNER_INVALID_ID"},
		{"RunnerInternalError", RunnerInternalError, "RUNNER_INTERNAL_ERROR"},
		{"RunnerGitLabOAuthFailed", RunnerGitLabOAuthFailed, "RUNNER_GITLAB_OAUTH_FAILED"},
		{"RunnerGitLabTokenInvalid", RunnerGitLabTokenInvalid, "RUNNER_GITLAB_TOKEN_INVALID"},
		{"RunnerGitLabInstanceUnknown", RunnerGitLabInstanceUnknown, "RUNNER_GITLAB_INSTANCE_UNKNOWN"},
		{"RunnerGitLabWebhookFailed", RunnerGitLabWebhookFailed, "RUNNER_GITLAB_WEBHOOK_FAILED"},
		{"RunnerGitLabScopeDenied", RunnerGitLabScopeDenied, "RUNNER_GITLAB_SCOPE_DENIED"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
