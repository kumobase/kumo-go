package types

import "testing"

func TestRDSRoundTrip(t *testing.T) {
	roundTrip(t, "CreateRDSInstanceRequest", CreateRDSInstanceRequest{
		Name:          "my-pg",
		Engine:        RDSEnginePostgreSQL,
		EngineVersion: "16",
		Plan:          "kumo.pg.small",
		StorageGB:     20,
	})
	sz := 50
	roundTrip(t, "UpdateRDSInstanceRequest+storage", UpdateRDSInstanceRequest{
		StorageGB: &sz,
	})
	roundTrip(t, "UpdateRDSInstanceRequest+plan", UpdateRDSInstanceRequest{
		Plan: "kumo.pg.medium",
	})
	roundTrip(t, "PublicRDSPlanResponse", PublicRDSPlanResponse{
		Slug:      "kumo.pg.small",
		Engine:    RDSEnginePostgreSQL,
		Name:      "Small",
		CPUvCPU:   "1",
		MemoryMB:  2048,
		PriceHour: "12.5000",
	})
	roundTrip(t, "PublicRDSPlanResponse+availability", PublicRDSPlanResponse{
		Slug:         "kumo.pg.large",
		Engine:       RDSEnginePostgreSQL,
		Name:         "Large",
		PriceHour:    "50.0000",
		Availability: &Availability{Available: false, Reason: AvailabilityReasonMemoryFull},
	})
	roundTrip(t, "RDSInstanceResponse", RDSInstanceResponse{
		ID:            7,
		Name:          "my-pg",
		Engine:        RDSEnginePostgreSQL,
		EngineVersion: "16",
		Mode:          string(RDSModeStandalone),
		Replicas:      1,
		Plan:          &PublicRDSPlanResponse{Slug: "kumo.pg.small", Engine: RDSEnginePostgreSQL, Name: "Small", PriceHour: "12.5000"},
		StorageGB:     20,
		Status:        string(RDSStatusSuspended),
		EndpointHost:  "my-pg.db.kumo.example",
		EndpointPort:  5432,
		IsSuspended:   true,
		SuspendReason: "insufficient balance",
	})
	roundTrip(t, "RDSMutationResponse", RDSMutationResponse{
		ID:          7,
		OperationID: "f1e2d3c4",
		Status:      string(RDSStatusProvisioning),
	})
	roundTrip(t, "RDSOperationResponse", RDSOperationResponse{
		OperationID: "f1e2d3c4",
		ActionType:  "create",
		Status:      string(RDSOperationInProgress),
	})
	roundTrip(t, "RDSConnectionResponse", RDSConnectionResponse{
		Host:     "my-pg.db.kumo.example",
		Port:     5432,
		Username: "kumo",
		Database: "kumo",
		Password: "s3cr3t",
		SSLMode:  "require",
	})
}
