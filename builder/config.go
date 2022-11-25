package builder

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	GoVer     string   `json:"gover"`
	Module    string   `json:"module"`
	Commit    string   `json:"commit"`
	Packages  []string `json:"packages"`
	WareHouse string   `json:"warehouse"`
	NetRC     string   `json:"netrc"`
}

var _defaultConfig = Config{
	GoVer:     "1.18",
	WareHouse: "/opt/go-dynamic-warehouse",
}

// ParseRemote parses a remote string into struct
// example: github.com/aura-studio/dynamic/builder@af3e5e21
func ParseRemote(remote string, packages ...string) []Config {
	strs := strings.Split(remote, "@")
	mod := strs[0]
	commit := strs[1]
	config := Config{
		GoVer:     _defaultConfig.GoVer,
		Module:    mod,
		Commit:    commit,
		Packages:  packages,
		WareHouse: _defaultConfig.WareHouse,
		NetRC:     "",
	}
	return []Config{config}
}

// BuildFromRemote builds a package from remote
func BuildFromRemote(remote string, packages ...string) {
	configs := ParseRemote(remote, packages...)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// ParseJson parses a json string into struct
func ParseJSON(str string) []Config {
	var configs []Config
	err := json.Unmarshal([]byte(str), &configs)
	if err != nil {
		log.Panic(err)
	}
	for i := range configs {
		if configs[i].GoVer == "" {
			configs[i].GoVer = _defaultConfig.GoVer
		}
		if configs[i].WareHouse == "" {
			configs[i].WareHouse = _defaultConfig.WareHouse
		}
	}
	return configs
}

// BuildFromJSON builds a package from json string
func BuildFromJSON(str string) {
	configs := ParseJSON(str)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// ParseJSONFile parses a json file into struct
func ParseJSONFile(path string) []Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	return ParseJSON(string(data))
}

// BuildFromJSONFile builds a package from json file
func BuildFromJSONFile(path string) {
	configs := ParseJSONFile(path)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// ParseJSONPath parses a json path into struct
func ParseJSONPath(path string) []Config {
	data, err := os.ReadFile(filepath.Join(path, "dynamic.json"))
	if err != nil {
		log.Panic(err)
	}
	return ParseJSON(string(data))
}

// BuildFromJSONPath builds a package from json path
func BuildFromJSONPath(path string) {
	configs := ParseJSONPath(path)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

type RenderData struct {
	Name      string
	Version   string
	Package   string
	Module    string
	House     string
	GoVersion string
	NetRC     string
}

func (c *Config) ToRenderData() []*RenderData {
	if len(c.Packages) == 0 {
		name := c.Module[strings.LastIndex(c.Module, "/")+1:]
		renderData := &RenderData{
			GoVersion: c.GoVer,
			Name:      name,
			Version:   c.Commit,
			Package:   c.Module,
			Module:    c.Module,
			House:     c.WareHouse,
			NetRC:     c.NetRC,
		}
		return []*RenderData{renderData}
	}

	renderDatas := make([]*RenderData, len(c.Packages))
	for i, pkg := range c.Packages {
		name := pkg[strings.LastIndex(pkg, "/")+1:]
		renderData := &RenderData{
			GoVersion: c.GoVer,
			Name:      name,
			Version:   c.Commit,
			Package:   strings.Join([]string{c.Module, pkg}, "/"),
			Module:    c.Module,
			House:     c.WareHouse,
			NetRC:     c.NetRC,
		}
		renderDatas[i] = renderData
	}

	return renderDatas
}
