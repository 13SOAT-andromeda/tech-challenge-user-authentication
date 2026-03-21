package domain

import "tech-challenge-user-validation/pkg/encryption"

type Password struct {
	hashed string
	hasher encryption.Hasher
}

func NewPasswordFromHash(hashed string, hasher encryption.Hasher) Password {
	return Password{
		hashed: hashed,
		hasher: hasher,
	}
}

func (p Password) GetHashed() string {
	return p.hashed
}

func (p Password) Compare(password string) error {
	return p.hasher.Compare(p.hashed, password)
}
