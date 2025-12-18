package dynamic

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
)

type Local struct {
	localPath string
}

func NewLocal(localPath string) *Local {
	if localPath == "" {
		return nil
	}

	return &Local{localPath: localPath}
}

func (l Local) Path() string {
	if l.localPath != "" {
		return l.localPath
	} else if runtime.GOOS == "windows" {
		return "C:/warehouse"
	} else {
		return "/opt/warehouse"
	}
}

func (l Local) Exists(name string) bool {
	localCgoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libcgo_%s.so", name))
	localGoFilePath := filepath.Join(l.Path(), toolchain.String(), name, fmt.Sprintf("libgo_%s.so", name))

	if stat, err := os.Stat(localCgoFilePath); err != nil || stat.Size() == 0 {
		log.Println("dynamic: Local Exists missing cgo file", localCgoFilePath)
		return false
	}

	if stat, err := os.Stat(localGoFilePath); err != nil || stat.Size() == 0 {
		log.Println("dynamic: Local Exists missing go file", localGoFilePath)
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
