package auth

// Event types and routing keys published by this service.
// Consumers may rely on these constants — bumping a value is a breaking change.
const (
	EventVersion = "v1"

	EventUserRegistered            = "user.registered"
	EventUserUpdated               = "user.updated"
	EventUserDeleted               = "user.deleted"
	EventUserPasswordResetRequested = "user.password_reset_requested"
	EventUserLogin                 = "user.login"
)

type UserRegisteredData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type UserUpdatedData struct {
	UserID         string   `json:"user_id"`
	ChangedFields  []string `json:"changed_fields"`
}

type UserDeletedData struct {
	UserID string `json:"user_id"`
}

type UserPasswordResetRequestedData struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
	ExpiresAt  string `json:"expires_at"`
}

type UserLoginData struct {
	UserID    string `json:"user_id"`
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
