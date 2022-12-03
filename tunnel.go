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
	warehouse string
	mu        sync.Mutex
	tunnelMap = make(map[string]Tunnel)
	remote    Remote
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

	if remote != nil {
		remote.MustExists(name)
	}

	var (
		plug *plugin.Plugin
		err  error
	)
	plug, err = plugin.Open(fmt.Sprintf("%s/%s/libgo_%s.so", warehouse, name, name))
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

func WithRemote(r Remote) {
	mu.Lock()
	defer mu.Unlock()

	remote = r
}
