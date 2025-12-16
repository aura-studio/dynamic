package dynamic

import (
	"regexp"
	"strings"
	"sync"
)

const (
	Default = "default"
	Latest  = "latest"
)

var (
	namespace string
	pkgMap    = make(map[string]Tunnel)
	muPkg     sync.Mutex
	allowed   = regexp.MustCompile(`^[A-Za-z0-9.-]+$`)
)

func UseNamespace(s string) {
	if !allowed.MatchString(s) {
		panic("invalid namespace")
	}
	if s == "" {
		s = Default
	}
	namespace = s
}

func GetPackage(pkg string, version string) (Tunnel, error) {
	muPkg.Lock()
	defer muPkg.Unlock()

	return getPackage(pkg, version)
}

func getPackage(pkg string, version string) (Tunnel, error) {
	name := getPackageTunnelName(pkg, version)

	if tunnel, ok := pkgMap[name]; ok {
		if tunnel != nil {
			return tunnel, nil
		} else if version != Latest {
			return getPackage(pkg, Latest)
		} else {
			return nil, ErrTunnelNotExits
		}
	} else {
		tunnel, err := GetTunnel(name)
		if err != nil {
			if !isTunnelNotExist(err) {
				return nil, err
			}

			pkgMap[name] = nil

			if version != Latest {
				return getPackage(pkg, Latest)
			}
			return nil, nil
		}
		pkgMap[name] = tunnel
		return tunnel, nil
	}
}

func ClosePackage(pkg string, version string) {
	muPkg.Lock()
	defer muPkg.Unlock()

	name := getPackageTunnelName(pkg, version)

	if tunnel, ok := pkgMap[name]; ok {
		if tunnel != nil {
			tunnel.Close()
		}
		delete(pkgMap, name)
	}
}

func RegisterPackage(pkg string, version string, tunnel Tunnel) {
	muPkg.Lock()
	defer muPkg.Unlock()

	if !allowed.MatchString(pkg) || !allowed.MatchString(version) {
		panic("invalid package name or version")
	}

	name := getPackageTunnelName(pkg, version)
	pkgMap[name] = tunnel
	RegisterTunnel(name, tunnel)
}

func getPackageTunnelName(pkg string, version string) string {
	return strings.Join([]string{namespace, pkg, version}, "_")
}
