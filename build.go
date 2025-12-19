package dynamic

// These values can be injected at build time, e.g.:
//
//	go build -ldflags "-X dynamic.DynamicOS=windows -X dynamic.DynamicArch=amd64 -X dynamic.BuildCompiler=gc -X dynamic.DynamicVariant=prod" ./...
var (
	DynamicOS       string
	DynamicArch     string
	DynamicCompiler string
	DynamicVariant  string
)
