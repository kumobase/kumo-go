package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/kumobase/kumo-go/client"
	"github.com/kumobase/kumo-go/codes"
	"github.com/kumobase/kumo-go/types"
)

func TestAuth_Refresh(t *testing.T) {
	var gotBody types.RefreshRequest
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/auth/refresh" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		writeStruct(w, 200, "", "token refreshed", &types.RefreshResponse{
			Token:        "new.access.jwt",
			RefreshToken: "kumo_rt_newrotatedtoken",
		})
	})

	out, err := c.Auth().Refresh(context.Background(), &types.RefreshRequest{
		RefreshToken: "kumo_rt_oldtoken",
		RememberMe:   true,
	})
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if gotBody.RefreshToken != "kumo_rt_oldtoken" || !gotBody.RememberMe {
		t.Fatalf("server saw wrong body: %+v", gotBody)
	}
	if out.Token != "new.access.jwt" || out.RefreshToken != "kumo_rt_newrotatedtoken" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestAuth_Refresh_Reused(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeStruct(w, 401, codes.RefreshTokenReused, "session revoked", nil)
	})

	_, err := c.Auth().Refresh(context.Background(), &types.RefreshRequest{RefreshToken: "kumo_rt_stale"})
	if err == nil {
		t.Fatal("expected error for reused token")
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != codes.RefreshTokenReused {
		t.Fatalf("expected APIError with code %s, got %v", codes.RefreshTokenReused, err)
	}
}

func TestAuth_Logout_Idempotent(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/logout" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		writeStruct(w, 200, "", "logged out", nil)
	})
	if err := c.Auth().Logout(context.Background(), &types.RefreshRequest{}); err != nil {
		t.Fatalf("Logout: %v", err)
	}
}
