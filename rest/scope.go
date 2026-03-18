package rest

// Scope represents a GitHub authorization scope.
//
// See https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps
type Scope string

const (
	ScopeRepo           Scope = "repo"
	ScopeRepoStatus     Scope = "repo:status"
	ScopeRepoDeployment Scope = "repo_deployment"
	ScopePublicRepo     Scope = "public_repo"
	ScopeRepoInvite     Scope = "repo:invite"
	ScopeSecurityEvents Scope = "security_events"
	ScopeAdminRepoHook  Scope = "admin:repo_hook"
	ScopeWriteRepoHook  Scope = "write:repo_hook"
	ScopeReadRepoHook   Scope = "read:repo_hook"
	ScopeAdminOrg       Scope = "admin:org"
	ScopeWriteOrg       Scope = "write:org"
	ScopeReadOrg        Scope = "read:org"
	ScopeAdminPublicKey Scope = "admin:public_key"
	ScopeWritePublicKey Scope = "write:public_key"
	ScopeReadPublicKey  Scope = "read:public_key"
	ScopeAdminOrgHook   Scope = "admin:org_hook"
	ScopeGist           Scope = "gist"
	ScopeNotifications  Scope = "notifications"
	ScopeUser           Scope = "user"
	ScopeReadUser       Scope = "read:user"
	ScopeUserEmail      Scope = "user:email"
	ScopeUserFollow     Scope = "user:follow"
	ScopeProject        Scope = "project"
	ScopeReadProject    Scope = "read:project"
	ScopeDeleteRepo     Scope = "delete_repo"
	ScopeWritePackages  Scope = "write:packages"
	ScopeReadPackages   Scope = "read:packages"
	ScopeDeletePackages Scope = "delete:packages"
	ScopeAdminGPGKey    Scope = "admin:gpg_key"
	ScopeWriteGPGKey    Scope = "write:gpg_key"
	ScopeReadGPGKey     Scope = "read:gpg_key"
	ScopeCodespace      Scope = "codespace"
	ScopeWorkflow       Scope = "workflow"
	ScopeReadAuditLog   Scope = "audit_log"
)
