package domain

import "testing"

func TestToken(t *testing.T) {
	token := Token{
		TokenID: "some-jti-uuid",
		UserID:  123,
	}

	if token.TokenID != "some-jti-uuid" {
		t.Errorf("expected TokenID some-jti-uuid, got %s", token.TokenID)
	}
	if token.UserID != 123 {
		t.Errorf("expected UserID 123, got %d", token.UserID)
	}
}
