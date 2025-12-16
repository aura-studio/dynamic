package dynamic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
)

func getEnv() string {
	os := GetOS()
	arch := GetArch()
	compiler := GetCompiler()
	variant := getVariant()
	env := os + "_" + arch + "_" + compiler + "_" + variant
	return env
}

// GetOS returns a best-effort OS descriptor.
// - On Linux with /etc/os-release: "<id><version_id>" like "ubuntu22.04".
// - On Windows (best-effort): "windows<version>" like "windows10.0.22631".
// - On Darwin/macOS (best-effort): "darwin<version>" like "darwin14.2.1".
// - Fallback: returns GOOS like "linux"/"darwin"/"windows".
func GetOS() string {
	goos := strings.ToLower(strings.TrimSpace(getGoEnvVar("GOOS")))
	if goos == "" {
		goos = strings.ToLower(strings.TrimSpace(runtime.GOOS))
	}

	switch goos {
	case "linux":
		if desc := detectLinuxDescriptor(); desc != "" {
			return desc
		}
		return "linux"
	case "windows":
		if v := detectWindowsVersion(); v != "" {
			return "windows" + v
		}
		return "windows"
	case "darwin":
		if v := detectDarwinVersion(); v != "" {
			return "darwin" + v
		}
		return "darwin"
	default:
		return goos
	}
}

func getGoEnvVar(key string) string {
	out, err := exec.Command("go", "env", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func detectLinuxDescriptor() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return ""
	}
	defer f.Close()

	m := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		m[k] = strings.Trim(v, `"`)
	}
	id := strings.ToLower(strings.TrimSpace(m["ID"]))
	ver := strings.TrimSpace(m["VERSION_ID"])
	if id == "" || ver == "" {
		return ""
	}
	return id + strings.TrimSpace(ver)
}

func detectWindowsVersion() string {
	// Prefer PowerShell because parsing is easier and locale independent.
	out, err := exec.Command(
		"powershell",
		"-NoProfile",
		"-NonInteractive",
		"-Command",
		"[System.Environment]::OSVersion.Version.ToString()",
	).Output()
	if err == nil {
		v := strings.TrimSpace(string(out))
		if v != "" {
			return v
		}
	}

	// Fallback to `cmd /c ver` (may vary by locale; extract first version-like token).
	out, err = exec.Command("cmd", "/c", "ver").Output()
	if err != nil {
		return ""
	}
	return extractFirstVersionLikeToken(string(out))
}

func detectDarwinVersion() string {
	// Prefer macOS product version: e.g. 14.2.1
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err == nil {
		v := strings.TrimSpace(string(out))
		if v != "" {
			return v
		}
	}
	// Fallback to kernel release: e.g. 23.2.0
	out, err = exec.Command("uname", "-r").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func extractFirstVersionLikeToken(s string) string {
	s = strings.TrimSpace(s)
	start := -1
	for i, r := range s {
		if unicode.IsDigit(r) {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	end := start
	for end < len(s) {
		c := s[end]
		if (c >= '0' && c <= '9') || c == '.' {
			end++
			continue
		}
		break
	}
	if end <= start {
		return ""
	}
	return s[start:end]
}

// GetArch returns current arch with best-effort variant.
// Examples: amd64, armv7, arm64v8.
func GetArch() string {
	// Prefer `go env -json` so we only spawn `go` once and can assemble
	// arch variants consistently (e.g. amd64v1, armv7, arm64v8).
	env, ok := getGoEnvArchVars()
	if ok {
		goarch := strings.ToLower(strings.TrimSpace(env.GOARCH))
		switch goarch {
		case "":
			// fall through to runtime below
		case "amd64":
			goamd64 := strings.ToLower(strings.TrimSpace(env.GOAMD64))
			// GOAMD64 is typically "v1"/"v2"/"v3"/"v4".
			if strings.HasPrefix(goamd64, "v") {
				return "amd64" + goamd64
			} else {
				fmt.Printf("unexpected GOAMD64 value: %q", goamd64)
			}
			return "amd64"
		case "arm":
			goarm := strings.TrimSpace(env.GOARM)
			if goarm != "" {
				return "armv" + goarm
			} else {
				fmt.Println("unexpected GOARM empty value")
			}
			return "arm"
		case "arm64":
			// Go does not expose an arm64 variant; treat as v8 by default.
			return "arm64v8"
		default:
			return goarch
		}
	}

	// Fallback when `go` is unavailable.
	goarch := strings.ToLower(strings.TrimSpace(runtime.GOARCH))
	if goarch == "arm64" {
		return "arm64v8"
	}
	return goarch
}

type goEnvArchVars struct {
	GOARCH  string `json:"GOARCH"`
	GOAMD64 string `json:"GOAMD64"`
	GOARM   string `json:"GOARM"`
}

func getGoEnvArchVars() (goEnvArchVars, bool) {
	out, err := exec.Command("go", "env", "-json").Output()
	if err != nil {
		return goEnvArchVars{}, false
	}
	var v goEnvArchVars
	if err := json.Unmarshal(out, &v); err != nil {
		return goEnvArchVars{}, false
	}
	return v, true
}

// GetCompiler returns `go env GOVERSION`, e.g. go1.20.5.
func GetCompiler() string {
	out, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		return strings.TrimSpace(runtime.Version())
	}
	return strings.TrimSpace(string(out))
}
