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
	if t.OS == "" {
		t.OS = env.GetOS()
	}

	if BuildArch != "" {
		t.Arch = BuildArch
	} else {
		t.Arch = os.Getenv("DYNAMIC_ARCH")
	}
	if t.Arch == "" {
		t.Arch = env.GetArch()
	}

	if BuileCompiler != "" {
		t.Compiler = BuileCompiler
	} else {
		t.Compiler = os.Getenv("DYNAMIC_COMPILER")
	}
	if t.Compiler == "" {
		t.Compiler = env.GetCompiler()
	}

	if BuildVariant != "" {
		t.Variant = BuildVariant
	} else {
		t.Variant = os.Getenv("DYNAMIC_VARIANT")
	}
	if t.Variant == "" {
		t.Variant = "generic" // 包含构建参数和so的路径都必须固定
	}
}

func (t Toolchain) String() string {
	return t.OS + "_" + t.Arch + "_" + t.Compiler + "_" + t.Variant
}
