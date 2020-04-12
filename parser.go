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
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
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
	mod, err := goListModules(mv.Path)
	if err != nil {
		return err
	}

	if mod.Indirect {
		// Skip indirect dependency
		return nil
	}

	if mod.Dir == "" {
		err := goModDownload(mv.Path)
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
			err := goModDownload(mu.Path + "@" + mu.Version)
			if err != nil {
				return err
			}

			uMod, err := goListModules(mu.Path + "@" + mu.Version)
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

func goModDownload(pkgId string) error {
	_, stdErr, err := goCmd("mod", "download", "-json", pkgId)
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stdErr)
	}
	return nil
}

func goListModules(pkgId string) (*Module, error) {
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
