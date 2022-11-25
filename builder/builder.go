package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

var templateMap = map[string]string{}

type RenderData struct {
	Name      string
	Version   string
	Package   string
	Module    string
	House     string
	GoVersion string
	NetRC     string
}

type Builder struct {
	config       *RenderData
	netRCPath    string
	netRCBakPath string
}

var config = &RenderData{
	Name:      "system",
	Version:   "e6e3d5fa",                                                               // git commit id (short, 8 chars)
	Package:   "codeup.aliyun.com/5efaebedaf80fad018f122fb/RogueLite/xCutServer/system", // package name (rel from Module)
	Module:    "codeup.aliyun.com/5efaebedaf80fad018f122fb/RogueLite/xCutServer",
	House:     "/opt/go-dynamic-warehouse",
	GoVersion: "1.18",
	NetRC:     `machine codeup.aliyun.com login finance@mirroringtech.com password agp_12d68a9572c2897119934df1d07456ec`,
}

var builder = NewBuilder(config)

func NewBuilder(c *RenderData) *Builder {
	user, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return &Builder{
		config:       c,
		netRCPath:    filepath.Join(user.HomeDir, ".netrc"),
		netRCBakPath: filepath.Join(user.HomeDir, ".netrc.go_dynamic_bak"),
	}
}

func (b *Builder) build() {
	fmt.Println("start...")
	defer fmt.Println("done!")

	b.bakNetRC()
	b.writeNetRc()
	defer b.restoreNetRC()

	b.generate()
	b.runBuilder()
}

// bakNetRC backup netrc file if exsits
func (b *Builder) bakNetRC() {
	fmt.Println("bakup", b.netRCPath)
	if _, err := os.Stat(b.netRCPath); err == nil {
		if err := os.Rename(b.netRCPath, b.netRCBakPath); err != nil {
			log.Panic(err)
		}
	}
}

// writeNetRc write netrc file from Builder.config.NetRC
func (b *Builder) writeNetRc() {
	fmt.Println("write", b.netRCPath)
	f, err := os.Create(b.netRCPath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(b.config.NetRC); err != nil {
		log.Panic(err)
	}
}

// restoreNetRC restore netrc file if exsits
func (b *Builder) restoreNetRC() {
	fmt.Println("restore", b.netRCPath)
	if _, err := os.Stat(b.netRCPath); err == nil {
		if err := os.Rename(b.netRCBakPath, b.netRCPath); err != nil {
			log.Panic(err)
		}
	}
}

// generate generate packages, go.mod and builder.sh
func (b *Builder) generate() {
	for pathTemplateStr, textTemplateStr := range templateMap {
		var pathBuilder strings.Builder
		if pathTemplate, err := template.New("dynamic").Parse(pathTemplateStr); err != nil {
			log.Panic(err)
		} else if err := pathTemplate.Execute(&pathBuilder, b.config); err != nil {
			log.Panic(err)
		}

		path := pathBuilder.String()
		fmt.Println("generate", path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Panic(err)
		}

		f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		if textTemplate, err := template.New("dynamic").Parse(textTemplateStr); err != nil {
			log.Panic(err)
		} else if err := textTemplate.Execute(f, b.config); err != nil {
			log.Panic(err)
		}
	}
}

// runBuilder run ./builder.sh
func (b *Builder) runBuilder() {
	cmd := exec.Command("bash", "-c", "./builder.sh")
	cmd.Dir = filepath.Join(b.config.House, b.config.Name+"_"+b.config.Version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// cat builder.sh
	builderPath := filepath.Join(cmd.Dir, "builder.sh")
	builderFile, err := os.Open(builderPath)
	if err != nil {
		log.Panic(err)
	}
	builderContent, err := io.ReadAll(builderFile)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("run builder.sh\n================\n%s================\n", string(builderContent))
	builderFile.Close()

	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}
}
