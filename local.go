package dynamic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

type Local struct {
	localPath string
}

func NewLocal() *Local {
	return &Local{}
}

func (l Local) Path() string {
	switch toolchain.Variant {
	case "generic":
		return "/opt/warehouse"
	}
	panic("dynamic: unsupported toolchain variant: " + toolchain.Variant)
}

func (l Local) Exists(name string) bool {
	localCgoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libcgo_%s.so", name))
	localGoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libgo_%s.so", name))

	// libgo is required.
	if stat, err := os.Stat(localGoFilePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	} else if stat.Size() == 0 {
		return false
	}

	// libcgo is optional: some builds may not produce it.
	if stat, err := os.Stat(localCgoFilePath); err != nil {
		if os.IsNotExist(err) {
			return true
		}
		return false
	} else if stat.Size() == 0 {
		return false
	}

	return true
}

func (l Local) PluginLoad(name string) (any, error) {
	localGoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libgo_%s.so", name))
	plug, err := plugin.Open(localGoFilePath)
	if err != nil {
		return nil, err
	}

	if symbol, err := plug.Lookup("Tunnel"); err == nil {
		return symbol, nil
	} else if symbol, err = plug.Lookup("New"); err == nil {
		newFunc, ok := symbol.(func() any)
		if !ok {
			return nil, errors.New("dynamic: unexpected type from symbol New")
		}
		return newFunc(), nil
	} else {
		return nil, err
	}
}
