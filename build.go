package dynamic

// These values can be injected at build time, e.g.:
//
//	go build -ldflags "-X dynamic.BuildOS=windows -X dynamic.BuildArch=amd64 -X dynamic.BuildCompiler=gc -X dynamic.BuildVariant=prod" ./...
var (
	BuildOS       string
	BuildArch     string
	BuileCompiler string
	BuildVariant  string
)
