package domain

type Token struct {
	TokenID string `json:"token_id"` // jti
	UserID  int64  `json:"user_id"`
}
