package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/colorstring"
	"github.com/radeksimko/go-mod-stat/go-src/cmd/go/_internal/modfile"
	"github.com/radeksimko/go-mod-stat/go-src/cmd/go/_internal/module"
)

type Parser struct {
	OutputWriter io.Writer
}

func (p *Parser) ParseModfile(path string) error {
	// Parse go modules file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	f, err := modfile.Parse(path, data, nil)
	if err != nil {
		return err
	}

	for _, r := range f.Require {
		err := p.parseModuleVersionRequirement(r.Mod)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseModuleVersionRequirement(mv module.Version) error {
	mod, err := getModuleData(mv.Path, "")
	if err != nil {
		return err
	}

	if mod.Indirect {
		// Skip indirect dependency
		return nil
	}

	if mod.Dir == "" {
		_, _, err := goCmd("mod", "download", "-json", mv.Path)
		if err != nil {
			return err
		}
		// Retry after downloading missing module
		return p.parseModuleVersionRequirement(mv)
	}

	if _, err = os.Stat(filepath.Join(mod.Dir, "go.mod")); os.IsNotExist(err) {
		colorstring.Fprintf(p.OutputWriter, "%s @ %s is [bold][red]module-unaware[reset]", mod.Path, mod.Version)

		if mod.Update != nil {
			// Check go.mod in latest version if update is available
			mu := mod.Update
			_, stdErr, err := goCmd("mod", "download", "-json", mu.Path+"@"+mu.Version)
			if err != nil {
				return fmt.Errorf("%s\n%s", err, stdErr)
			}

			uMod, err := getModuleData(mu.Path, mu.Version)
			if err != nil {
				return err
			}

			if _, err = os.Stat(filepath.Join(uMod.Dir, "go.mod")); err == nil {
				colorstring.Fprintf(p.OutputWriter, " [bold][yellow](updatable to %s)[reset]", uMod.Version)
			}
		}
		fmt.Println("")
	}

	return nil
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
		return nil, stderr.String(), fmt.Errorf("%q: %s", args, err)
	}

	return &stdout, stderr.String(), nil
}
