package dynamic

import (
	"os"
)

func init() {
	toolchain.Init()
}

type Toolchain struct {
	OS       string
	Arch     string
	Compiler string
	Variant  string
}

var toolchain = NewToolchain()

func NewToolchain() *Toolchain {
	return &Toolchain{}
}

func (t *Toolchain) Init() {
	if BuildOS != "" {
		t.OS = BuildOS
	} else {
		t.OS = os.Getenv("DYNAMIC_OS")
	}

	if BuildArch != "" {
		t.Arch = BuildArch
	} else {
		t.Arch = os.Getenv("DYNAMIC_ARCH")
	}

	if BuileCompiler != "" {
		t.Compiler = BuileCompiler
	} else {
		t.Compiler = os.Getenv("DYNAMIC_COMPILER")
	}

	if BuildVariant != "" {
		t.Variant = BuildVariant
	} else {
		t.Variant = os.Getenv("DYNAMIC_VARIANT")
	}

	if t.OS == "" || t.Arch == "" || t.Compiler == "" || t.Variant == "" {
		panic("dynamic: OS, Arch, Compiler, Variant must be set via -ldflags or environment variables")
	}
}

func (t Toolchain) String() string {
	return t.OS + "_" + t.Arch + "_" + t.Compiler + "_" + t.Variant
}
