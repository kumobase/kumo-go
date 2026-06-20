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
	roundTrip(t, "CreateRDSInstanceRequest+tls", CreateRDSInstanceRequest{
		Name:          "my-pg",
		Engine:        RDSEnginePostgreSQL,
		EngineVersion: "16",
		Plan:          "kumo.pg.small",
		StorageGB:     20,
		TLSMode:       string(RDSTLSModeRequired),
	})
	roundTrip(t, "UpdateRDSTLSRequest", UpdateRDSTLSRequest{
		TLSMode: string(RDSTLSModeOptional),
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
	roundTrip(t, "PublicRDSEngineVersionResponse", PublicRDSEngineVersionResponse{
		Engine:    RDSEnginePostgreSQL,
		Version:   "16",
		IsDefault: true,
		Status:    "supported",
	})
	roundTrip(t, "RDSPgParameterResponse", RDSPgParameterResponse{
		Name: "max_connections", DataType: "int", ApplyMethod: "static",
		DefaultValue: "100", MinValue: "1", MaxValue: "8000", Unit: "", Description: "max client connections",
	})
	roundTrip(t, "PublicRDSParameterTemplateResponse", PublicRDSParameterTemplateResponse{
		ID: 3, Slug: "pg16-oltp", Name: "OLTP tuned", Engine: RDSEnginePostgreSQL, EngineVersion: "16",
		IsSystem: false, IsDefault: false,
		Parameters: []RDSParameter{{Name: "max_connections", Value: "200"}, {Name: "work_mem", Value: "8MB"}},
	})
	desc := "tuned"
	roundTrip(t, "UpdateRDSParameterTemplateRequest", UpdateRDSParameterTemplateRequest{
		Description: &desc,
		Parameters:  []RDSParameter{{Name: "work_mem", Value: "16MB"}},
	})
	roundTrip(t, "CreateRDSInstanceRequest+template", CreateRDSInstanceRequest{
		Name: "my-pg", Engine: RDSEnginePostgreSQL, EngineVersion: "16", Plan: "kumo.pg.small",
		StorageGB: 20, ParameterTemplate: "pg16-oltp",
	})
	rr := 2
	roundTrip(t, "CreateRDSInstanceRequest+ha", CreateRDSInstanceRequest{
		Name: "my-pg", Engine: RDSEnginePostgreSQL, EngineVersion: "16", Plan: "kumo.pg.small",
		StorageGB: 20, Mode: string(RDSModeHA), ReadReplicas: &rr,
	})
	roundTrip(t, "UpdateRDSInstanceRequest+replicas", UpdateRDSInstanceRequest{
		Mode: string(RDSModeHA), ReadReplicas: &rr,
	})
	roundTrip(t, "CreateRDSInstanceRequest+replica-specs", CreateRDSInstanceRequest{
		Name: "my-pg", Engine: RDSEnginePostgreSQL, EngineVersion: "16", Plan: "kumo.pg.medium",
		StorageGB: 20, ReadReplicas: &rr,
		ReadReplicaSpecs: []ReadReplicaSpec{{Plan: "kumo.pg.small"}, {Plan: "kumo.pg.small"}},
	})
	roundTrip(t, "UpdateRDSInstanceRequest+replica-specs", UpdateRDSInstanceRequest{
		ReadReplicas:     &rr,
		ReadReplicaSpecs: []ReadReplicaSpec{{Plan: "kumo.pg.small"}, {Plan: "kumo.pg.large"}},
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
		TLSMode:       string(RDSTLSModeRequired),
	})
	roundTrip(t, "RDSInstanceResponse+heterogeneous-replicas", RDSInstanceResponse{
		ID:            8,
		Name:          "my-pg-ha",
		Engine:        RDSEnginePostgreSQL,
		EngineVersion: "16",
		Mode:          string(RDSModeHA),
		Replicas:      4, // 1 primary + 1 sync standby + 2 read replicas
		ReadReplicas:  2,
		ReadReplicaDetails: []ReadReplicaDetail{
			{Ordinal: 0, Plan: &PublicRDSPlanResponse{Slug: "kumo.pg.small", Engine: RDSEnginePostgreSQL, Name: "Small", PriceHour: "12.5000"}, Status: string(RDSStatusReady)},
			{Ordinal: 1, Plan: &PublicRDSPlanResponse{Slug: "kumo.pg.large", Engine: RDSEnginePostgreSQL, Name: "Large", PriceHour: "50.0000"}, Status: string(RDSStatusReady)},
		},
		Plan:      &PublicRDSPlanResponse{Slug: "kumo.pg.medium", Engine: RDSEnginePostgreSQL, Name: "Medium", PriceHour: "25.0000"},
		StorageGB: 20,
		Status:    string(RDSStatusReady),
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
		Host:             "my-pg.db.kumo.example",
		ReadHost:         "my-pg-ro.db.kumo.example",
		ReadReplicaHosts: []string{"my-pg-0.db.kumo.example", "my-pg-2.db.kumo.example"},
		Port:             5432,
		Username:         "kumo",
		Database:         "kumo",
		Password:         "s3cr3t",
		SSLMode:          "require",
		CACert:           "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----\n",
	})
	roundTrip(t, "RDSConnectionResponse+nossl", RDSConnectionResponse{
		Host:     "my-pg.db.kumo.example",
		Port:     5432,
		Username: "kumo",
		Database: "kumo",
		Password: "s3cr3t",
		SSLMode:  "disable",
	})
}
