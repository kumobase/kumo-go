package types

import "time"

// SecretType enumerates the kinds of secret a user can create. The wire
// value is what the server stores and what's returned in every secret
// response's Type field.
type SecretType string

const (
	SecretTypeRegistry    SecretType = "registry"
	SecretTypeEnvVar      SecretType = "env_var"
	SecretTypeFile        SecretType = "file"
	SecretTypeCertificate SecretType = "certificate"
)

// SecretRegistry holds Docker-registry pull credentials for SecretTypeRegistry
// secrets. Username and Password are 3..100 chars server-side; RegistryHost
// is optional (defaults to Docker Hub when empty).
type SecretRegistry struct {
	RegistryHost string `json:"registry_host,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

// CertificateContent holds the cert + key pair for SecretTypeCertificate
// secrets. Both fields are required when the parent secret's Type is
// "certificate".
type CertificateContent struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

// SecretConsumer is one app that references a secret. Populated in
// ResponseGetSecret.UsedBy so callers can show "in use by N apps" before a
// delete attempt (which would otherwise return 409 SECRET_IN_USE).
//
// UsageType is "attached" when the secret is mounted via secret_var_apps,
// or "tls" when it's directly attached as the app's TLS secret.
type SecretConsumer struct {
	AppID              uint   `json:"app_id"`
	AppName            string `json:"app_name"`
	UsageType          string `json:"usage_type"`
	SourceFrom         string `json:"source_from,omitempty"`
	MountTo            string `json:"mount_to,omitempty"`
	RestartWhenUpdated bool   `json:"restart_when_updated"`
}

// RequestSecretBase is the shared head of create/update payloads.
type RequestSecretBase struct {
	Name string     `json:"name"` // 3..100 chars
	Type SecretType `json:"type"`
}

// CreateSecretRequest is the body for POST /api/v1/secrets. Honors
// Idempotency-Key. Exactly one of SecretRegistry / EnvironmentVariables /
// FileContent / CertificateContent is required, matching Type.
type CreateSecretRequest struct {
	RequestSecretBase
	SecretRegistry       SecretRegistry        `json:"secret_registry,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables,omitempty"`
	FileContent          string                `json:"file_content,omitempty"`
	CertificateContent   CertificateContent    `json:"certificate_content,omitempty"`
}

// UpdateSecretRequest is the body for PATCH /api/v1/secrets/:id. The Type
// field is immutable — the server returns 400 SECRET_TYPE_IMMUTABLE on a
// change attempt. Supports If-Match for optimistic concurrency.
type UpdateSecretRequest struct {
	CreateSecretRequest
}

// GetSecretAllResponse is the list-item shape returned by GET /api/v1/secrets.
// The sensitive payload (registry password, env values, file content, key)
// is intentionally NOT included in list responses — fetch the secret by ID
// to retrieve it.
type GetSecretAllResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Type        SecretType `json:"type"`
	UsedByCount int64      `json:"used_by_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ResponseGetSecret is the detail shape returned by GET /api/v1/secrets/:id.
// Server sets ETag from UpdatedAt; echo it back in If-Match on PATCH.
//
// Sensitive payload is included here (the GET-detail endpoint is the only
// way to read a secret's contents back). Only one of the four payload
// fields will be populated, matching Type.
type ResponseGetSecret struct {
	ID                   uint                  `json:"id"`
	Name                 string                `json:"name"`
	Type                 SecretType            `json:"type"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	UsedByCount          int64                 `json:"used_by_count"`
	UsedBy               []SecretConsumer      `json:"used_by"`
	SecretRegistry       *SecretRegistry       `json:"secret_registry,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables,omitempty"`
	FileContent          string                `json:"file_content,omitempty"`
	CertificateContent   *CertificateContent   `json:"certificate_content,omitempty"`
}
