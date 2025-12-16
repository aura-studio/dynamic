package dynamic

import "runtime"

func getEnv() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	compiler := runtime.Version()
	variant := getVariant()
	env := os + "_" + arch + "_" + compiler + "_" + variant
	return env
}

func getVariant() string {
	// TODO: 根据debuginfo，推断variable的变种
	return "plain"
}
