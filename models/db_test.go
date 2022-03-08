package models

import (
	"testing"
)

func TestInit(t *testing.T) {
	_, err := Init()

	if err != nil {
		t.Errorf("When loading DB, got error: %s", err)
	}

}