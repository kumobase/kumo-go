package codes

import "testing"

func TestVoucherCampaignCodeValues(t *testing.T) {
	cases := map[string]string{
		VoucherCampaignNotFound:                    "VOUCHER_CAMPAIGN_NOT_FOUND",
		VoucherCampaignAlreadyExists:               "VOUCHER_CAMPAIGN_ALREADY_EXISTS",
		VoucherCampaignInvalidState:                "VOUCHER_CAMPAIGN_INVALID_STATE",
		VoucherCampaignRulesLocked:                 "VOUCHER_CAMPAIGN_RULES_LOCKED",
		VoucherCampaignInvalidRule:                 "VOUCHER_CAMPAIGN_INVALID_RULE",
		VoucherCampaignInvalidScope:                "VOUCHER_CAMPAIGN_INVALID_SCOPE",
		VoucherCampaignPoolInsufficient:            "VOUCHER_CAMPAIGN_POOL_INSUFFICIENT",
		VoucherCampaignApplicationNotFound:         "VOUCHER_CAMPAIGN_APPLICATION_NOT_FOUND",
		VoucherCampaignReversalInsufficientBalance: "VOUCHER_CAMPAIGN_REVERSAL_INSUFFICIENT_BALANCE",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("code = %q, want %q", got, want)
		}
	}
}
