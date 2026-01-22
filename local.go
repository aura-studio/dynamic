package dynamic

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

type Local struct {
	localPath string
}

func NewLocal(localPath string) *Local {
	return &Local{
		localPath: localPath,
	}
}

func (l Local) Path() string {
	return l.localPath
}

func (l Local) Exists(name string) bool {
	localCgoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libcgo_%s.so", name))
	localGoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libgo_%s.so", name))
	log.Printf("[dynamic] check warehouse package %s go file: %s", name, localGoFilePath)
	log.Printf("[dynamic] check warehouse package %s cgo file: %s", name, localCgoFilePath)

	// libgo is required.
	if stat, err := os.Stat(localGoFilePath); err != nil {
		log.Printf("[dynamic] stat error: %v", err)
		return false
	} else if stat.Size() == 0 {
		log.Printf("[dynamic] file size is zero: %s", localGoFilePath)
		return false
	}

	// libcgo is required.
	if stat, err := os.Stat(localCgoFilePath); err != nil {
		log.Printf("[dynamic] stat error: %v", err)
		return false
	} else if stat.Size() == 0 {
		log.Printf("[dynamic] file size is zero: %s", localCgoFilePath)
		return false
	}

	log.Printf("[dynamic] found warehouse package %s go file: %s", name, localGoFilePath)
	log.Printf("[dynamic] found warehouse package %s cgo file: %s", name, localCgoFilePath)
	return true
}

func (l Local) Load(name string) (any, error) {
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
