package dynamic

func UseWarehouse(local, remote string) {
	if !allowed.IsPath(local) {
		panic("dynamic: invalid local warehouse path")
	}
	if !allowed.IsURL(remote) {
		panic("dynamic: invalid remote warehouse URL")
	}
	warehouse.Init(local, remote)
}

func UseNamespace(namespace string) {
	if !allowed.IsKeyword(namespace) {
		panic("dynamic: invalid package namespace")
	}
	packageCenter.UseNamespace(namespace)
}

func UseDefaultVersion(version string) {
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid default package version")
	}
	packageCenter.UseDefaultVersion(version)
}

func RegisterPackage(packageName string, version string, tunnel Tunnel) {
	if !allowed.IsKeyword(packageName) {
		panic("dynamic: invalid package name")
	}
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid package version")
	}
	packageCenter.RegisterPackage(packageName, version, tunnel)
}

func GetPackage(packageName string, version string) (Tunnel, error) {
	if !allowed.IsKeyword(packageName) {
		panic("dynamic: invalid package name")
	}
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid package version")
	}
	tunnel, err := packageCenter.GetTunnel(packageName, version)
	if err != nil {
		return nil, err
	}
	return tunnel, nil
}

func ClosePackage(packageName string, version string) {
	if !allowed.IsKeyword(packageName) {
		panic("dynamic: invalid package name")
	}
	if !allowed.IsKeyword(version) {
		panic("dynamic: invalid package version")
	}
	packageCenter.ClosePackage(packageName, version)
}
