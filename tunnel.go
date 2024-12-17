package dynamic

import (
	"fmt"
	"path/filepath"
	"plugin"
	"runtime"
	"sync"
)

type Tunnel interface {
	Init()
	Invoke(string, string) string
	Close()
}

type Template struct {
}

func (t *Template) Init() {
}

func (t *Template) Close() {
}

func (t *Template) Invoke(name string, args string) string {
	return ""
}

var (
	mu        sync.Mutex
	tunnelMap = make(map[string]Tunnel)
)

func GetWarehouse() string {
	if local != "" {
		return local
	} else if runtime.GOOS == "windows" {
		return "C:/warehouse"
	} else {
		return "/tmp/warehouse"
	}
}

func GetTunnel(name string) (Tunnel, error) {
	mu.Lock()
	defer mu.Unlock()

	if tunnel, ok := tunnelMap[name]; ok {
		return tunnel, nil
	}

	remote := NewRemote()
	if err := remote.Sync(name); err != nil {
		return nil, err
	}

	var (
		plug *plugin.Plugin
		err  error
	)
	localFileName := fmt.Sprintf("libgo_%s.so", name)
	localFilePath := filepath.Join(GetWarehouse(), runtime.Version(), name, localFileName)
	plug, err = plugin.Open(localFilePath)
	if err != nil {
		return nil, err
	}

	var (
		tunnel Tunnel
		ok     bool
	)

	if symbol, err := plug.Lookup("Tunnel"); err == nil {
		tunnel, ok = symbol.(Tunnel)
		if !ok {
			return nil, fmt.Errorf("unexpected type from symbol Tunnel: %s", name)
		}
	} else if symbol, err = plug.Lookup("New"); err == nil {
		newFunc, ok := symbol.(func() Tunnel)
		if !ok {
			return nil, fmt.Errorf("unexpected type from symbol New: %s", name)
		}
		tunnel = newFunc()
	} else {
		return nil, err
	}

	tunnel.Init()

	tunnelMap[name] = tunnel

	return tunnel, nil
}

func CloseTunnel(name string) error {
	mu.Lock()
	defer mu.Unlock()

	if tunnel, ok := tunnelMap[name]; ok {
		tunnel.Close()
		delete(tunnelMap, name)
	}

	return nil
}

func RangeTunnel(f func(string, Tunnel) bool) {
	mu.Lock()
	defer mu.Unlock()

	for name, tunnel := range tunnelMap {
		if !f(name, tunnel) {
			break
		}
	}
}

// RegisterTunnel is usually used in debug mode
func RegisterTunnel(name string, tunnel Tunnel) {
	mu.Lock()
	defer mu.Unlock()

	tunnel.Init()
	tunnelMap[name] = tunnel
}
