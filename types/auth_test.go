package types

import "testing"

func TestAuthRefreshRoundTrip(t *testing.T) {
	roundTrip(t, "RefreshRequest/cookie", RefreshRequest{})
	roundTrip(t, "RefreshRequest/body", RefreshRequest{
		RefreshToken: "kumo_rt_deadbeefdeadbeefdeadbeefdeadbeef",
		RememberMe:   true,
	})
	roundTrip(t, "RefreshResponse", RefreshResponse{
		Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.sig",
		RefreshToken: "kumo_rt_cafebabecafebabecafebabecafebabe",
	})
}
