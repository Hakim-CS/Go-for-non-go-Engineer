package servitor

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewServitorAlpha(t *testing.T) {
	sa := NewServitorAlpha()

	if sa == nil {
		t.Errorf("NewServitoAlpha() return nil  ")
	}
}

func TestDefaultPassword(t *testing.T) {
	sa := NewServitorAlpha()
	defaultPassword, err := sa.DefaultPassword()

	if err != nil {
		t.Errorf("DefaultPassword() returned an error: %v", err)
	}

	if _, err := uuid.Parse(defaultPassword); err != nil {
		t.Errorf("DefaultPassword() returned an invalied uuid : %v ", err)
	}

}
