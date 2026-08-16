package types

import "time"

type VoucherCampaignTrigger string

const (
	VoucherCampaignTriggerTopupSuccess  VoucherCampaignTrigger = "topup_success"
	VoucherCampaignTriggerBillingCharge VoucherCampaignTrigger = "billing_charge"
)

type VoucherCampaignBenefitType string

const (
	VoucherCampaignBenefitPercentageBonus    VoucherCampaignBenefitType = "percentage_bonus"
	VoucherCampaignBenefitFixedBonus         VoucherCampaignBenefitType = "fixed_bonus"
	VoucherCampaignBenefitPercentageDiscount VoucherCampaignBenefitType = "percentage_discount"
	VoucherCampaignBenefitFixedDiscount      VoucherCampaignBenefitType = "fixed_discount"
)

type VoucherCampaignAudience string

const (
	VoucherCampaignAudienceAll                  VoucherCampaignAudience = "all"
	VoucherCampaignAudienceNeverSuccessfulTopup VoucherCampaignAudience = "never_successful_topup"
)

type VoucherCampaignStatus string

const (
	VoucherCampaignStatusDraft     VoucherCampaignStatus = "draft"
	VoucherCampaignStatusScheduled VoucherCampaignStatus = "scheduled"
	VoucherCampaignStatusActive    VoucherCampaignStatus = "active"
	VoucherCampaignStatusPaused    VoucherCampaignStatus = "paused"
	VoucherCampaignStatusEnded     VoucherCampaignStatus = "ended"
)

type VoucherCampaignAvailability string

const (
	VoucherCampaignAvailabilityAvailable        VoucherCampaignAvailability = "available"
	VoucherCampaignAvailabilityInactive         VoucherCampaignAvailability = "inactive"
	VoucherCampaignAvailabilityIneligible       VoucherCampaignAvailability = "ineligible"
	VoucherCampaignAvailabilityUserLimitReached VoucherCampaignAvailability = "user_limit_reached"
	VoucherCampaignAvailabilityBudgetExhausted  VoucherCampaignAvailability = "budget_exhausted"
	VoucherCampaignAvailabilityPoolInsufficient VoucherCampaignAvailability = "pool_insufficient"
)

type VoucherApplicationStatus string

const (
	VoucherApplicationPending      VoucherApplicationStatus = "pending"
	VoucherApplicationNotAwarded   VoucherApplicationStatus = "not_awarded"
	VoucherApplicationReserved     VoucherApplicationStatus = "reserved"
	VoucherApplicationApplied      VoucherApplicationStatus = "applied"
	VoucherApplicationReversed     VoucherApplicationStatus = "reversed"
	VoucherApplicationManualReview VoucherApplicationStatus = "manual_review"
	VoucherApplicationWaived       VoucherApplicationStatus = "waived"
)

type VoucherApplicationResolutionAction string

const (
	VoucherApplicationResolutionRetry VoucherApplicationResolutionAction = "retry"
	VoucherApplicationResolutionWaive VoucherApplicationResolutionAction = "waive"
)

type VoucherCampaignProductType string

const (
	VoucherCampaignProductApp               VoucherCampaignProductType = "app"
	VoucherCampaignProductJobs              VoucherCampaignProductType = "jobs"
	VoucherCampaignProductStorage           VoucherCampaignProductType = "storage"
	VoucherCampaignProductContainerRegistry VoucherCampaignProductType = "container_registry"
)

type VoucherCampaignScope struct {
	ProductType VoucherCampaignProductType `json:"product_type"`
	PlanID      *uint                      `json:"plan_id,omitempty"`
}

type CreateVoucherCampaignRequest struct {
	Slug         string                     `json:"slug"`
	Name         string                     `json:"name"`
	Description  string                     `json:"description,omitempty"`
	Trigger      VoucherCampaignTrigger     `json:"trigger"`
	BenefitType  VoucherCampaignBenefitType `json:"benefit_type"`
	Audience     VoucherCampaignAudience    `json:"audience"`
	Currency     string                     `json:"currency"`
	Value        string                     `json:"value"`
	PerEventCap  *string                    `json:"per_event_cap,omitempty"`
	PerUserLimit *int                       `json:"per_user_limit,omitempty"`
	Budget       *string                    `json:"budget,omitempty"`
	Priority     int                        `json:"priority"`
	Timezone     string                     `json:"timezone"`
	StartsAt     time.Time                  `json:"starts_at"`
	EndsAt       time.Time                  `json:"ends_at"`
	Scopes       []VoucherCampaignScope     `json:"scopes,omitempty"`
}

// Trigger, BenefitType, and Currency are deliberately absent: those wire
// rules are immutable after creation, including while the campaign is draft.
type UpdateVoucherCampaignRequest struct {
	Slug           *string                  `json:"slug,omitempty"`
	Name           *string                  `json:"name,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	Audience       *VoucherCampaignAudience `json:"audience,omitempty"`
	Value          *string                  `json:"value,omitempty"`
	PerEventCap    *string                  `json:"per_event_cap,omitempty"`
	ClearEventCap  bool                     `json:"clear_event_cap,omitempty"`
	PerUserLimit   *int                     `json:"per_user_limit,omitempty"`
	ClearUserLimit bool                     `json:"clear_user_limit,omitempty"`
	Budget         *string                  `json:"budget,omitempty"`
	ClearBudget    bool                     `json:"clear_budget,omitempty"`
	Priority       *int                     `json:"priority,omitempty"`
	Timezone       *string                  `json:"timezone,omitempty"`
	StartsAt       *time.Time               `json:"starts_at,omitempty"`
	EndsAt         *time.Time               `json:"ends_at,omitempty"`
	Scopes         *[]VoucherCampaignScope  `json:"scopes,omitempty"`
}

type VoucherCampaignActionRequest struct {
	Reason string `json:"reason"`
}

type VoucherApplicationActionRequest struct {
	Reason string `json:"reason"`
}

type VoucherCampaignResponse struct {
	ID              uint                       `json:"id"`
	Slug            string                     `json:"slug"`
	Name            string                     `json:"name"`
	Description     string                     `json:"description,omitempty"`
	Trigger         VoucherCampaignTrigger     `json:"trigger"`
	BenefitType     VoucherCampaignBenefitType `json:"benefit_type"`
	Audience        VoucherCampaignAudience    `json:"audience"`
	Currency        string                     `json:"currency"`
	Value           string                     `json:"value"`
	PerEventCap     *string                    `json:"per_event_cap,omitempty"`
	PerUserLimit    *int                       `json:"per_user_limit,omitempty"`
	Budget          *string                    `json:"budget,omitempty"`
	Priority        int                        `json:"priority"`
	Timezone        string                     `json:"timezone"`
	StartsAt        time.Time                  `json:"starts_at"`
	EndsAt          time.Time                  `json:"ends_at"`
	Status          VoucherCampaignStatus      `json:"status"`
	Scopes          []VoucherCampaignScope     `json:"scopes"`
	ReservedAmount  string                     `json:"reserved_amount"`
	AppliedAmount   string                     `json:"applied_amount"`
	AvailableAmount *string                    `json:"available_amount,omitempty"`
	ScheduledAt     *time.Time                 `json:"scheduled_at,omitempty"`
	PausedAt        *time.Time                 `json:"paused_at,omitempty"`
	EndedAt         *time.Time                 `json:"ended_at,omitempty"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
}

type PublicVoucherCampaignResponse struct {
	ID                  uint                        `json:"id"`
	Slug                string                      `json:"slug"`
	Name                string                      `json:"name"`
	Description         string                      `json:"description,omitempty"`
	Trigger             VoucherCampaignTrigger      `json:"trigger"`
	BenefitType         VoucherCampaignBenefitType  `json:"benefit_type"`
	Currency            string                      `json:"currency"`
	Value               string                      `json:"value"`
	PerEventCap         *string                     `json:"per_event_cap,omitempty"`
	Timezone            string                      `json:"timezone"`
	StartsAt            time.Time                   `json:"starts_at"`
	EndsAt              time.Time                   `json:"ends_at"`
	Scopes              []VoucherCampaignScope      `json:"scopes"`
	Eligible            bool                        `json:"eligible"`
	Availability        VoucherCampaignAvailability `json:"availability"`
	IneligibilityReason string                      `json:"ineligibility_reason,omitempty"`
}

type VoucherApplicationResponse struct {
	ID                  uint                                `json:"id"`
	CampaignID          uint                                `json:"campaign_id"`
	CampaignSlug        string                              `json:"campaign_slug"`
	CampaignName        string                              `json:"campaign_name"`
	Status              VoucherApplicationStatus            `json:"status"`
	ResolutionAction    *VoucherApplicationResolutionAction `json:"resolution_action,omitempty"`
	SourceType          string                              `json:"source_type"`
	SourceReference     string                              `json:"source_reference"`
	GrossAmount         string                              `json:"gross_amount"`
	BenefitAmount       string                              `json:"benefit_amount"`
	NetAmount           string                              `json:"net_amount"`
	Currency            string                              `json:"currency"`
	LedgerTransactionID *string                             `json:"ledger_transaction_id,omitempty"`
	FailureReason       string                              `json:"failure_reason,omitempty"`
	AppliedAt           *time.Time                          `json:"applied_at,omitempty"`
	ReversedAt          *time.Time                          `json:"reversed_at,omitempty"`
	CreatedAt           time.Time                           `json:"created_at"`
	UpdatedAt           time.Time                           `json:"updated_at"`
}
