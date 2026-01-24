package dynamic

import (
	"errors"
	"log"
	"strings"
	"sync"
)

const (
	NamespaceDefault = "default"
	VersionDefault   = "default"
	VersionLatest    = "latest"
)

type DynamicIndex struct {
	Namespace string
	Package   string
	Version   string
}

func NewDynamicIndex(namespace, pkg, version string) *DynamicIndex {
	return &DynamicIndex{
		Namespace: namespace,
		Package:   pkg,
		Version:   version,
	}
}

func (d DynamicIndex) String() string {
	return strings.Join([]string{d.Namespace, d.Package, d.Version}, "_")
}

type Dynamic struct {
	index  DynamicIndex
	tunnel Tunnel
}

func NewDynamic(index DynamicIndex, tunnel Tunnel) *Dynamic {
	return &Dynamic{
		index:  index,
		tunnel: tunnel,
	}
}

func (d *Dynamic) GetTunnel() Tunnel {
	return d.tunnel
}

type DynamicCenter struct {
	namesapce      string
	defaultVersion string
	mu             sync.Mutex
	dynamics       map[DynamicIndex]*Dynamic
}

var packageCenter = NewPackageCenter()

func NewPackageCenter() *DynamicCenter {
	return &DynamicCenter{
		namesapce:      NamespaceDefault,
		defaultVersion: VersionDefault,
		dynamics:       make(map[DynamicIndex]*Dynamic),
	}
}

func (dc *DynamicCenter) UseNamespace(s string) {
	dc.namesapce = s
}

func (dc *DynamicCenter) UseDefaultVersion(v string) {
	dc.defaultVersion = v
}

func (dc *DynamicCenter) GetTunnel(pkg string, version string) (tunnel Tunnel, err error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	var index DynamicIndex

	// first try with provided version
	index = *NewDynamicIndex(dc.namesapce, pkg, version)

	if dynamic, ok := dc.dynamics[index]; ok {
		return dynamic.GetTunnel(), nil
	}

	if tunnel, err := tunnelCenter.GetTunnel(index.String()); err == nil {
		dc.cache(pkg, version, tunnel)
		return tunnel, nil
	} else {
		log.Printf("[dynamic] get tunnel %s failed: %v", index.String(), err)
	}

	// then try with default version
	index = *NewDynamicIndex(dc.namesapce, pkg, dc.defaultVersion)

	if dynamic, ok := dc.dynamics[index]; ok {
		dc.cache(pkg, version, dynamic.GetTunnel())
		return dynamic.GetTunnel(), nil
	}

	if tunnel, err := tunnelCenter.GetTunnel(index.String()); err == nil {
		dc.cache(pkg, version, tunnel)
		dc.cache(pkg, dc.defaultVersion, tunnel)
		return tunnel, nil
	} else {
		log.Printf("[dynamic] get tunnel %s failed: %v", index.String(), err)
	}

	return nil, errors.New("dynamic: both provided version and default version not found")
}

func (dc *DynamicCenter) ClosePackage(pkg string, version string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	index := *NewDynamicIndex(dc.namesapce, pkg, version)
	if dynamic, ok := dc.dynamics[index]; ok {
		if dynamic != nil {
			dynamic.GetTunnel().Close()
		}
		delete(dc.dynamics, index)
	}
}

func (dc *DynamicCenter) RegisterPackage(pkg string, version string, tunnel Tunnel) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	index := dc.cache(pkg, version, tunnel)
	tunnelCenter.RegisterTunnel(index.String(), tunnel)
}

func (dc *DynamicCenter) cache(pkg string, version string, tunnel Tunnel) DynamicIndex {
	index := *NewDynamicIndex(dc.namesapce, pkg, version)
	dyanmic := NewDynamic(index, tunnel)
	dc.dynamics[index] = dyanmic
	return index
}
