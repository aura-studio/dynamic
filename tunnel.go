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
	warehouse = "/opt/go-dynamic-warehouse"
	mu        sync.Mutex
	tunnelMap = make(map[string]Tunnel)
)

func init() {
	if env, ok := os.LookupEnv("GO_DYNAMIC_WAREHOUSE"); ok {
		warehouse = env
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

	tunnelMap[name] = tunnel

	return tunnel, nil
}
