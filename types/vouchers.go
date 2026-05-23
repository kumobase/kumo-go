package types

import "time"

// Voucher SDK surface is read-only — POST /api/v1/vouchers/redeem is
// sessionOnly (financial mutation) and is not callable with an API key.
// The redeem request type is omitted for that reason; only the history
// item shape is exposed.

// RedemptionHistoryResponse is one row of GET /api/v1/vouchers/history.
// Amount is a decimal string (IDR).
type RedemptionHistoryResponse struct {
	ID          uint      `json:"id"`
	VoucherCode string    `json:"voucher_code"`
	Amount      string    `json:"amount"`
	RedeemedAt  time.Time `json:"redeemed_at"`
}
