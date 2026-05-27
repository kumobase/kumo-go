package types

import (
	"testing"
	"time"
)

// One round-trip per module's primary DTOs. We don't enumerate every field —
// the round-trip itself catches tag drift; this just exercises the shape so
// CI fails fast if a JSON tag silently changes.

func TestApps_RoundTrip(t *testing.T) {
	now := time.Date(2026, 5, 23, 10, 0, 0, 0, time.UTC)
	roundTrip(t, "CreateAppRequest", CreateAppRequest{
		BaseCreateApp: BaseCreateApp{
			Name: "my-app", Image: "nginx:1.27", Port: 8080,
			IsExposed: true, Replicas: 2,
			Autoscaling: &AutoscalingConfig{Enabled: true, MinReplicas: 1, MaxReplicas: 5, CPUTargetPercentage: intPtr(70)},
		},
		EnvironmentVariables: []EnvironmentVariable{{Key: "FOO", Value: "bar"}},
		PricingSlug:          "kumo.nano",
		SecretVars:           []SecretVar{{SecretId: 1, RestartWhenUpdated: true}},
		HealthCheck:          &HealthCheck{Type: "http", Path: "/health", Port: 8080},
	})
	roundTrip(t, "CreateAppResponse", CreateAppResponse{
		ID: 7, Name: "my-app", GenerateAppName: "my-app-x4z",
		DeploymentStatus: string(AppDeploymentStatusDeploying),
		OperationID:      "abc-123",
		UpdatedAt:        now,
	})
	roundTrip(t, "AppByIdResponse", AppByIdResponse{
		Id:                 1,
		CreateAppRequest:   CreateAppRequest{BaseCreateApp: BaseCreateApp{Name: "x", Image: "nginx", Port: 80, Replicas: 1}, PricingSlug: "kumo.nano"},
		GeneratedSubDomain: "x.kumo.run",
		Source:             AppSourceGitBuild,
		Language:           "static",
		OutputDir:          "dist",
		BuildCommand:       "build",
		AppStatus:          "running",
		StatusMessage:      "all good",
		DesiredReplicas:    1, ReadyReplicas: 1,
		CreatedAt: now, UpdatedAt: now,
	})
	roundTrip(t, "AddCustomDomainRequest", AddCustomDomainRequest{Domain: "shop.example.com"})
}

func TestVPS_RoundTrip(t *testing.T) {
	roundTrip(t, "RentServerRequest", RentServerRequest{Provider: "vultr", Region: "sgp", Plan: "vc2-1c-1gb", Name: "edge-1"})
	roundTrip(t, "PublicVPSPlanResponse", PublicVPSPlanResponse{
		ProviderName: "vultr", PlanID: 12, ExternalPlanID: "vc2-1c-1gb",
		Name: "1 CPU / 1GB", CPU: 1, Memory: 1024, Disk: 25,
		Egress: 1000, SellingPrice: "4.99",
		Availability: &Availability{Available: true},
	})
	roundTrip(t, "VPSServerResponse", VPSServerResponse{
		ID: 1, DisplayProvider: "vultr", RegionID: "sgp",
		Status: string(VPSStatusRunning), SSHPort: 22, AutoRenew: true,
		CreatedAt: "2026-05-23T10:00:00Z", SSHSetupCompleted: true,
	})
}

func TestSecrets_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "CreateSecretRequest", CreateSecretRequest{
		RequestSecretBase:    RequestSecretBase{Name: "db-creds", Type: SecretTypeEnvVar},
		EnvironmentVariables: []EnvironmentVariable{{Key: "DB_URL", Value: "postgres://x"}},
	})
	roundTrip(t, "ResponseGetSecret", ResponseGetSecret{
		ID: 1, Name: "tls", Type: SecretTypeCertificate,
		CreatedAt: now, UpdatedAt: now,
		CertificateContent: &CertificateContent{Certificate: "-----BEGIN-----", PrivateKey: "-----KEY-----"},
	})
}

func TestVolumes_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "CreateVolumeRequest", CreateVolumeRequest{
		Name: "data", StorageTier: "ssd-std", SizeGB: 10, MountPath: "/data",
	})
	roundTrip(t, "VolumeResponse", VolumeResponse{
		ID: 1, Name: "data", SizeGB: 10, MountPath: "/data",
		Status: string(VolumeStatusReady),
		StorageTier: StorageTierResponse{
			ID: 1, Slug: "ssd-std", Name: "Standard SSD",
			PricePerGBHour: "0.0001234", MinSizeGB: 1, MaxSizeGB: 1000,
		},
		CreatedAt: now, UpdatedAt: now,
	})
}

func TestRegistry_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "CreateRepositoryRequest", CreateRepositoryRequest{
		Name: "my-app", TagMutability: TagMutabilityImmutable, SoftDeleteDays: intPtr(7),
	})
	roundTrip(t, "RepositoryResponse", RepositoryResponse{
		ID: 1, Name: "my-app", TagMutability: TagMutabilityMutable,
		SoftDeleteDays: 7, CreatedAt: now, UpdatedAt: now,
	})
	roundTrip(t, "ManifestResponse", ManifestResponse{
		ID: 1, Digest: "sha256:abc", MediaType: "application/vnd.oci.image.manifest.v1+json",
		SizeBytes: 1024, ImageSizeBytes: 162203697, PushedAt: now,
	})
}

func TestOrganizations_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "CreateOrganizationRequest", CreateOrganizationRequest{
		Slug: "acme", DisplayName: "Acme Corp",
	})
	roundTrip(t, "OrganizationResponse", OrganizationResponse{
		ID: 1, Slug: "acme", DisplayName: "Acme Corp", OwnerUserID: 1,
		RegistryAutoCreateRepos: true, CreatedAt: now, UpdatedAt: now,
	})
}

func TestBilling_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "PublicChargeResponse", PublicChargeResponse{
		ID: 1, SubscriptionID: 7, ProductType: "vps", PlanName: "vultr-1c-1gb",
		Amount: "4.99", Currency: "IDR",
		PeriodStart: now, PeriodEnd: now, ChargeType: "prepaid", Status: "charged",
		ReferenceID: "ref-1", CreatedAt: now,
	})
	roundTrip(t, "BillingSummaryResponse", BillingSummaryResponse{
		Currency: "IDR",
		CurrentPeriod: PeriodSummary{
			Start: now, End: now, TotalCharged: "100.00",
			ByProduct: ProductBreakdown{VPS: "50.00", App: "30.00", Storage: "20.00"},
		},
		PreviousPeriodTotal: "200.00",
	})
}

func TestSmallModules_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "GetProfileResponse", GetProfileResponse{
		FullName: "Alice", Email: "alice@example.com",
		IsVerified: true, IsAdmin: false, HasPassword: true, HasGoogle: false,
	})
	roundTrip(t, "GetBalanceResponse", GetBalanceResponse{Balance: "1000.50"})
	roundTrip(t, "HasPasswordResponse", HasPasswordResponse{HasPassword: true})
	roundTrip(t, "RedemptionHistoryResponse", RedemptionHistoryResponse{
		ID: 1, VoucherCode: "ABCDEFGH", Amount: "100.00", RedeemedAt: now,
	})
	roundTrip(t, "ReferralStatsResponse", ReferralStatsResponse{
		Code: "ALICE-42", ReferralCount: 3, MaxReferrals: 10,
		TotalEarned: "300.00", PendingRewards: 1, Enabled: true,
	})
	roundTrip(t, "DashboardSummary", DashboardSummary{AppCount: 7})
	roundTrip(t, "TicketResponse", TicketResponse{
		ID: 1, DisplayID: "TKT-1", Subject: "help", Description: "broken",
		Category: "technical", Priority: "normal", Status: "open",
		CreatedAt: now, UpdatedAt: now,
	})
}

func TestAPIKeys_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "CreateAPIKeyRequest", CreateAPIKeyRequest{
		Name: "ci", ExpiresInDays: intPtr(30), Scopes: []string{"read", "write"},
	})
	roundTrip(t, "APIKeyResponse", APIKeyResponse{
		ID: 1, Name: "ci", KeyPrefix: "kumo_sk_abc",
		CreatedAt: now, Scopes: []string{"read", "write"},
	})
	roundTrip(t, "APIKeyCreateResponse", APIKeyCreateResponse{
		APIKeyResponse: APIKeyResponse{
			ID: 1, Name: "ci", KeyPrefix: "kumo_sk_abc",
			CreatedAt: now, Scopes: []string{"read"},
		},
		Key: "kumo_sk_abc_full_secret_value",
	})
}

func TestSourceConnect_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	roundTrip(t, "SourceConnectionResponse", SourceConnectionResponse{
		ID: 1, Provider: SourceProviderGitHub, InstallationID: 12345,
		AccountLogin: "acme", AccountType: "Organization",
		ManageURL: "https://github.com/organizations/acme/settings/installations/12345",
		Status:    SourceConnectionStatusActive, CreatedAt: now, UpdatedAt: now,
	})
	roundTrip(t, "SourceRepoResponse", SourceRepoResponse{
		ID: 99, FullName: "acme/web", Private: true, DefaultBranch: "main",
	})
}

func TestBuild_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	started := now.Add(time.Minute)
	finished := now.Add(5 * time.Minute)
	roundTrip(t, "BuildResponse", BuildResponse{
		ID: 1, AppID: 7, CommitSHA: "abc123", Ref: "refs/heads/main",
		Status: BuildStatusSucceeded, ImageDigest: "sha256:deadbeef",
		LogURL:    "https://logs.kumo.run/builds/1.txt?sig=x",
		CreatedAt: now, StartedAt: &started, FinishedAt: &finished,
	})
	roundTrip(t, "CreateGitBuildAppRequest", CreateGitBuildAppRequest{
		Name: "my-app", Port: 8080, IsExposed: true, Replicas: 2,
		RepoFullName: "acme/web", Branch: "main", Language: "static",
		OutputDir: "dist", BuildCommand: "build",
		EnvironmentVariables: []EnvironmentVariable{{Key: "FOO", Value: "bar"}},
		PricingSlug:          "kumo.nano",
		HealthCheck:          &HealthCheck{Type: "http", Path: "/health", Port: 8080},
	})
}

func intPtr(v int) *int { return &v }
