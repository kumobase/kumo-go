// Package version exposes the kumo-go SDK version string. Clients should
// use it in their HTTP User-Agent header so the server can surface compat
// warnings when an SDK falls outside the supported range.
package version

// SDKVersion is the SemVer of the kumo-go release this code was built
// against. The constant is bumped in the same commit that creates the git
// tag, so a build is always traceable to a published tag.
const SDKVersion = "v0.15.0"
