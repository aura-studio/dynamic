package dynamic

// Compatibility helpers used by remote implementation.
// Keep these thin wrappers so other components can evolve independently.

func getLocalWarehouse() string {
	if warehouse.Local == nil {
		return ""
	}
	return warehouse.Local.GetPath()
}

func getToolChain() string {
	return toolchain.String()
}
