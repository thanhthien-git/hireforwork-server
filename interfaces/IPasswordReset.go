package interfaces

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordReset struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}
