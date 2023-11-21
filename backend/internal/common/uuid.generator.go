package common

import (
	"github.com/gofrs/uuid"
)

type UUIDGenerator struct{}

func (u *UUIDGenerator) Gen() string {
	return uuid.Must(uuid.NewV4()).String()
}

func (u *UUIDGenerator) Valid(s string) bool {
	_, err := uuid.FromString(s)
	if err != nil {
		return false
	}
	return true
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}
