package dynamic

import (
	"errors"
	"log"
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

type TunnelCenter struct {
	mu      sync.Mutex
	tunnels map[string]Tunnel
}

var tunnelCenter = NewTunnelCenter()

func NewTunnelCenter() *TunnelCenter {
	return &TunnelCenter{
		tunnels: make(map[string]Tunnel),
	}
}

func (tc *TunnelCenter) GetTunnel(name string) (Tunnel, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	log.Println("dynamic: TunnelCenter GetTunnel", name)

	if tunnel, ok := tc.tunnels[name]; ok {
		return tunnel, nil
	}

	log.Printf("%#v", warehouse)

	plugin, err := warehouse.Load(name)
	if err != nil {
		return nil, err
	}

	tunnel, ok := plugin.(Tunnel)
	if !ok {
		return nil, errors.New("dynamic: symbol is not a Tunnel")
	}

	tunnel.Init()
	tc.tunnels[name] = tunnel

	return tunnel, nil
}

func (tc *TunnelCenter) CloseTunnel(name string) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tunnel, ok := tc.tunnels[name]; ok {
		tunnel.Close()
		delete(tc.tunnels, name)
	}

	return nil
}

func (tc *TunnelCenter) RangeTunnel(f func(string, Tunnel) bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	for name, tunnel := range tc.tunnels {
		if !f(name, tunnel) {
			break
		}
	}
}

// RegisterTunnel is usually used in debug mode
func (tc *TunnelCenter) RegisterTunnel(name string, tunnel Tunnel) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tunnel.Init()
	tc.tunnels[name] = tunnel
}
