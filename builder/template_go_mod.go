package main

func init() {
	templateMap["{{.House}}/{{.Name}}_{{.Version}}/go.mod"] = `module dynamicbuilder

go {{.GoVersion}}

require (
	{{.Module}} {{.Version}}
)
`
}
