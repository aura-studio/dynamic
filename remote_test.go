package dynamic_test

import (
	"testing"

	"github.com/aura-studio/dynamic"
)

func TestTunner(t *testing.T) {
	if _, err := dynamic.GetTunnel("testdynamic1_test"); err != nil {
		t.Error(err)
	}

	if _, err := dynamic.GetTunnel("testdynamic2_test"); err != nil {
		t.Error(err)
	}
}
