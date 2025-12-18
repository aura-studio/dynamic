package dynamic

import "log"

func UseWarehouse(local, remote string) {
	log.Printf("dynamic: UseWarehouse local_set=%q remote_set=%q", local, remote)
	if !allowed.IsPath(local) {
		panic("dynamic: invalid local warehouse path")
	}
	if !allowed.IsURL(remote) {
		panic("dynamic: invalid remote warehouse URL")
	}
	warehouse.Init(local, remote)
	log.Printf("%#v", warehouse)
}

func UseNamespace(namespace string) {
	log.Printf("dynamic: UseNamespace namespace=%q", namespace)
	if !allowed.IsKeyword(namespace) {
		panic("dynamic: invalid package namespace")
	}
	packageCenter.UseNamespace(namespace)
}

func UseDefaultVersion(version string) {
	log.Printf("dynamic: UseDefaultVersion version=%q", version)
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid default package version")
	}
	packageCenter.UseDefaultVersion(version)
}

func RegisterPackage(packageName string, version string, tunnel Tunnel) {
	log.Printf("dynamic: RegisterPackage package=%q version=%q", packageName, version)
	if !allowed.IsKeyword(packageName) {
		panic("dynamic: invalid package name")
	}
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid package version")
	}
	packageCenter.RegisterPackage(packageName, version, tunnel)
}

func GetPackage(packageName string, version string) (Tunnel, error) {
	log.Printf("dynamic: GetPackage package=%q version=%q", packageName, version)
	if !allowed.IsKeyword(packageName) {
		panic("dynamic: invalid package name")
	}
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid package version")
	}
	tunnel, err := packageCenter.GetTunnel(packageName, version)
	if err != nil {
		log.Printf("dynamic: GetPackage failed package=%q version=%q err=%v", packageName, version, err)
		return nil, err
	}
	log.Printf("dynamic: GetPackage ok package=%q version=%q", packageName, version)
	return tunnel, nil
}

func ClosePackage(packageName string, version string) {
	log.Printf("dynamic: ClosePackage package=%q version=%q", packageName, version)
	if !allowed.IsKeyword(packageName) {
		panic("dynamic: invalid package name")
	}
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid package version")
	}
	packageCenter.ClosePackage(packageName, version)
}
