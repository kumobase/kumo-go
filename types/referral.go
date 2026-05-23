package types

import "time"

// Referral SDK surface is read-only — the referral routes are sessionOnly
// today (account-level marketing data, not for machine credentials), so
// these shapes are exposed here for completeness in case the access policy
// relaxes in a future release.

// ReferralStatsResponse is returned by GET /api/v1/referral/stats.
// TotalEarned is a decimal string in IDR.
type ReferralStatsResponse struct {
	Code           string `json:"code"`
	ReferralCount  int64  `json:"referral_count"`
	MaxReferrals   int    `json:"max_referrals"`
	TotalEarned    string `json:"total_earned"`
	PendingRewards int64  `json:"pending_rewards"`
	Enabled        bool   `json:"enabled"`
}

// ReferralListItem is one row of GET /api/v1/referral/list.
// RewardAmount and TopUpAmount are decimal strings.
type ReferralListItem struct {
	ReferredEmail string    `json:"referred_email"`
	ReferredAt    time.Time `json:"referred_at"`
	RewardStatus  string    `json:"reward_status"`
	RewardAmount  string    `json:"reward_amount"`
	TopUpAmount   string    `json:"top_up_amount,omitempty"`
}
