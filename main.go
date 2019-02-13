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

		outBuffer, _, err := goCmd("list", "-json", "-m", m.Path)
		if err != nil {
			log.Fatal(err)
		}
		module := Module{}
		err = json.Unmarshal(outBuffer.Bytes(), &module)
		if err != nil {
			log.Fatal(err)
		}

		if !module.Indirect {
			if _, err = os.Stat(filepath.Join(module.Dir, "go.mod")); os.IsNotExist(err) {
				fmt.Printf("%s @ %s is not module-aware\n", module.Path, module.Version)
			}
		}
	}
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
