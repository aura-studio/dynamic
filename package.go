package dynamic

import (
	"strings"
	"sync"
)

const Latest = "latest"

var (
	namespace  string
	packageMap = make(map[string]Tunnel)
	muPackage  sync.Mutex
)

func UseNamespace(s string) {
	namespace = s
}

func GetPackage(packageName string, commit string) (Tunnel, error) {
	muPackage.Lock()
	defer muPackage.Unlock()

	return getPackage(packageName, commit)
}

func getPackage(packageName string, commit string) (Tunnel, error) {
	name := getPackageTunnelName(packageName, commit)

	if tunnel, ok := packageMap[name]; ok {
		return tunnel, nil
	} else {
		tunnel, err := GetTunnel(name)
		if err != nil {
			if !isTunnelNotExist(err) {
				return nil, err
			}

			packageMap[name] = nil

			if commit != Latest {
				return getPackage(packageName, Latest)
			}
			return nil, nil
		}
		packageMap[name] = tunnel
		return tunnel, nil
	}
}

func ClosePackage(packageName string, commit string) {
	muPackage.Lock()
	defer muPackage.Unlock()

	name := getPackageTunnelName(packageName, commit)

	if tunnel, ok := packageMap[name]; ok {
		if tunnel != nil {
			tunnel.Close()
		}
		delete(packageMap, name)
	}
}

func RegisterPackage(packageName string, commit string, tunnel Tunnel) {
	muPackage.Lock()
	defer muPackage.Unlock()

	name := getPackageTunnelName(packageName, commit)
	packageMap[name] = tunnel
	RegisterTunnel(name, tunnel)
}

func getPackageTunnelName(packageName string, commit string) string {
	name := strings.Join([]string{packageName, commit}, "_")
	if namespace != "" {
		name = strings.Join([]string{namespace, name}, "_")
	}
	return name
}
