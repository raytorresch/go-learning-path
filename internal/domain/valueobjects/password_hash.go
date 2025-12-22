package valueobjects

import "errors"

type PasswordHash struct {
	value string
}

func NewPasswordHash(value string) (pass PasswordHash, err error) {
	if value == "" {
		return PasswordHash{}, errors.New("password hash cannot be empty")
	}
	if len(value) < 8 {
		return PasswordHash{}, errors.New("invalid format")
	}
	return PasswordHash{value: value}, nil
}

func (p PasswordHash) Value() string {
	return p.value
}
