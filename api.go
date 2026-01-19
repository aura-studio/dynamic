package dynamic

// UseWarehouse:
//
//	如果使用此函数，local一定要有值，而remote可以为空。
//	local表示本地仓库路径，remote表示远程仓库URL。
//	remote如果有值，则会启用远程仓库同步功能。
//	Case 1: 不调用UseWarehouse函数，则不启用仓库功能，直走静态Package。
//	Case 2: 只调用UseWarehouse(local, ""), 则启用本地仓库功能，不启用远程同步功能。
//	Case 3: 调用UseWarehouse(local, remote), 则启用本地仓库功能，并启用远程同步功能。
func UseWarehouse(local, remote string) {
	if !allowed.IsPath(local) {
		panic("dynamic: invalid local warehouse path")
	}
	if remote != "" && !allowed.IsURL(remote) {
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
