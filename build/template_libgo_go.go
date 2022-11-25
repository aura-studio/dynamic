package main

func init() {
	templateMap["{{.House}}/{{.Name}}_{{.Version}}/libgo_{{.Name}}_{{.Version}}/libgo.go"] = `package main

/*
#cgo CFLAGS: -I{{.House}}/{{.Name}}_{{.Version}}/
#cgo LDFLAGS: -L{{.House}}/{{.Name}}_{{.Version}}/ -lcgo_{{.Name}}_{{.Version}}
#include "{{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}}.h"
#include "stdlib.h"
*/
import "C"
import "unsafe"

type tunnel struct{}

func (t tunnel) Init() {
	C.dynamic_cgo_{{.Name}}_{{.Version}}_init()
}

func (t tunnel) Invoke(route string, req string) string {
	rsp_cstr := C.dynamic_cgo_{{.Name}}_{{.Version}}_invoke(C.CString(route), C.CString(req))
	rsp := C.GoString(rsp_cstr)
	C.free(unsafe.Pointer(rsp_cstr))
	return rsp
}

func (t tunnel) Close() {
	C.dynamic_cgo_{{.Name}}_{{.Version}}_close()
}

var Tunnel tunnel

func main() {}
`
}
