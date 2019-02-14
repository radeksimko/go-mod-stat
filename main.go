package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/radeksimko/go-mod-stat/go-src/cmd/go/_internal/modfile"
)

const modRelPath = `pkg/mod`

type Module struct {
	Path      string       `json:",omitempty"` // module path
	Version   string       `json:",omitempty"` // module version
	Versions  []string     `json:",omitempty"` // available module versions
	Replace   *Module      `json:",omitempty"` // replaced by this module
	Time      *time.Time   `json:",omitempty"` // time version was created
	Update    *Module      `json:",omitempty"` // available update (with -u)
	Main      bool         `json:",omitempty"` // is this the main module?
	Indirect  bool         `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string       `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string       `json:",omitempty"` // path to go.mod file describing module, if any
	Error     *ModuleError `json:",omitempty"` // error loading module
	GoVersion string       `json:",omitempty"` // go version used in module
}

type ModuleError struct {
	Err string // error text
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Parse go modules file
	gomodPath := filepath.Join(cwd, "go.mod")
	data, err := ioutil.ReadFile(gomodPath)
	if err != nil {
		log.Fatal(err)
	}
	f, err := modfile.Parse(gomodPath, data, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range f.Require {
		m := r.Mod

		module, err := getModuleData(m.Path, "")
		if err != nil {
			log.Fatal(err)
		}

		if !module.Indirect {
			if module.Dir == "" {
				_, _, err := goCmd("mod", "download", "-json", m.Path)
				if err != nil {
					log.Fatal(err)
				}
			}

			if _, err = os.Stat(filepath.Join(module.Dir, "go.mod")); os.IsNotExist(err) {
				colorstring.Printf("%s @ %s is [bold][red]module-unaware[reset]", module.Path, module.Version)

				if module.Update != nil {
					// Check go.mod in latest version if update is available
					mu := module.Update
					_, _, err := goCmd("mod", "download", "-json", mu.Path+"@"+mu.Version)
					if err != nil {
						log.Fatal(err)
					}

					uModule, err := getModuleData(mu.Path, mu.Version)
					if err != nil {
						log.Fatal(err)
					}

					if _, err = os.Stat(filepath.Join(uModule.Dir, "go.mod")); err == nil {
						colorstring.Printf(" [bold][yellow](updatable to %s)[reset]", uModule.Version)
					}
				}

				fmt.Println("")
			}
		}
	}
}

func getModuleData(path, version string) (*Module, error) {
	pkgId := path
	if version != "" {
		pkgId += "@" + version
	}

	outBuffer, _, err := goCmd("list", "-json", "-u", "-m", pkgId)
	if err != nil {
		return nil, err
	}
	module := Module{}
	err = json.Unmarshal(outBuffer.Bytes(), &module)
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func goCmd(args ...string) (*bytes.Buffer, string, error) {
	cmd := exec.Command("go", args...)

	var stdout, stderr bytes.Buffer
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, stderr.String(), err
	}

	return &stdout, stderr.String(), nil
}
