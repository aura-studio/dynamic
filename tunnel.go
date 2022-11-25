package dynamic

import (
	"fmt"
	"os"
	"plugin"
	"sync"
)

type Tunnel interface {
	Init()
	Invoke(string, string) string
	Close()
}

var (
	warehouse string
	mu        sync.Mutex
	tunnelMap = make(map[string]Tunnel)
)

func init() {
	if env, ok := os.LookupEnv("GO_DYNAMIC_WAREHOUSE"); ok {
		warehouse = env
	} else {
		warehouse = "/opt/go-dynamic-warehouse"
	}
}

func GetTunnel(name string) (Tunnel, error) {
	mu.Lock()
	defer mu.Unlock()

	if tunnel, ok := tunnelMap[name]; ok {
		return tunnel, nil
	}

	var (
		plug *plugin.Plugin
		err  error
	)
	plug, err = plugin.Open(fmt.Sprintf("%s/%s/libgo_%s.so", warehouse, name, name))
	if err != nil {
		return nil, err
	}

	symbol, err := plug.Lookup("Tunnel")
	if err != nil {
		return nil, err
	}

	tunnel, ok := symbol.(Tunnel)
	if !ok {
		return nil, fmt.Errorf("unexpected type from symbol Tunnel: %s", name)
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

	tunnelMap[name] = tunnel
}
