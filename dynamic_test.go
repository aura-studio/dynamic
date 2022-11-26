package dynamic_test

import (
	"testing"

	"github.com/aura-studio/dynamic"
	. "github.com/frankban/quicktest"
)

// testv1 source
/*
package testmod

func Display() string {
	return "this is test mod v1.0.0"
}
*/

// testv2 source
/*
package testmod

import "fmt"

func Display() {
	fmt.Println("this is test mod v2.0.0")
}
*/

func TestDynamic(t *testing.T) {
	c := New(t)
	c.Run("TestCrossPackageVersion", func(c *C) {
		testdynamic1, err := dynamic.GetTunnel("testdynamic1_test")
		c.Assert(err, IsNil)
		testdynamic1.Invoke("", "")

		testdynamic2, err := dynamic.GetTunnel("testdynamic2_test")
		c.Assert(err, IsNil)
		testdynamic2.Invoke("", "")
	})
}
