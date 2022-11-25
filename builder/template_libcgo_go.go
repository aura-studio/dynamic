package main

func init() {
	templateMap["{{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}}/libcgo.go"] = `package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"{{.Package}}"
)

var tunnel = {{.Name}}.Tunnel

//export dynamic_cgo_{{.Name}}_{{.Version}}_init
func dynamic_cgo_{{.Name}}_{{.Version}}_init() {
	tunnel.Init()
}

//export dynamic_cgo_{{.Name}}_{{.Version}}_invoke
func dynamic_cgo_{{.Name}}_{{.Version}}_invoke(route_cstr *C.char, req_cstr *C.char) *C.char {
	route := C.GoString(route_cstr)
	C.free(unsafe.Pointer(route_cstr))

	req := C.GoString(req_cstr)
	C.free(unsafe.Pointer(req_cstr))

	return C.CString(tunnel.Invoke(route, req))
}

//export dynamic_cgo_{{.Name}}_{{.Version}}_close
func dynamic_cgo_{{.Name}}_{{.Version}}_close() {
	tunnel.Close()
}

func main() {}
`
}
