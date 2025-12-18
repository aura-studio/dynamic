package dynamic

import (
	"os"
)

// These values can be injected at build time, e.g.:
//
//	go build -ldflags "-X dynamic.OS=windows -X dynamic.Arch=amd64 -X dynamic.Compiler=gc -X dynamic.Variant=prod" ./...
//
// If not provided via -ldflags, init() will try to read them from environment variables:
//
//	DYNAMIC_OS, DYNAMIC_ARCH, DYNAMIC_COMPILER, DYNAMIC_VARIANT
var (
	OS       string
	Arch     string
	Compiler string
	Variant  string
)

func init() {
	if OS == "" {
		OS = os.Getenv("DYNAMIC_OS")
	}
	if Arch == "" {
		Arch = os.Getenv("DYNAMIC_ARCH")
	}
	if Compiler == "" {
		Compiler = os.Getenv("DYNAMIC_COMPILER")
	}
	if Variant == "" {
		Variant = os.Getenv("DYNAMIC_VARIANT")
	}
	if OS == "" || Arch == "" || Compiler == "" || Variant == "" {
		panic("dynamic: OS, Arch, Compiler, Variant must be set via -ldflags or environment variables")
	}
}

func getToolChain() string {
	return OS + "_" + Arch + "_" + Compiler + "_" + Variant
}
