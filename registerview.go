package registerviews

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/qor/qor/utils"
)

// DetectViewsDir 解决 go mod 模式无法注册 qor-admin 等包的 views
func DetectViewsDir(pkgorg, pkgname, subpath string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH

	}

	if abspath := isAbs(pkgname, subpath, pkgorg); len(abspath) > 0 {
		return abspath
	}

	ppath := filepath.Join(filepath.Join(gopath, "/pkg/mod/"), pkgorg)
	if _, err := os.Stat(ppath); err == nil {
		foundp, err := walkPath(ppath, pkgname, subpath)
		if err == nil {
			return foundp
		}
		color.Red(fmt.Sprintf("os.Stat2 error %v\n", err))
	} else {
		color.Red(fmt.Sprintf("os.Stat1 error %v\n", err))
	}

	return ""
}

func walkPath(ppath string, pkgname string, subpath string) (string, error) {
	var foundp string
	var found bool
	if err := filepath.Walk(ppath, func(p string, f os.FileInfo, err error) error { // nolint: errcheck, gosec, unparam
		if found {
			return nil
		}

		if foundp = foundPath(p, pkgname, subpath); len(foundp) > 0 {
			found = true
		}

		return nil

	}); err != nil {
		return "", err
	}

	return foundp, nil
}

func foundPath(p, pkgname, subpath string) string {
	if hasPerfix(p, pkgname) {
		if vp := filepath.Join(p, subpath+"views"); isExistingDir(vp) {
			return vp
		}
	}
	return ""
}

func hasPerfix(p string, pkgname string) bool {
	return strings.HasPrefix(filepath.Base(p), pkgname+"@")
}

func isAbs(pkgname string, subpath string, pkgorg string) string {
	if pkgname == "" && subpath == "" {
		if filepath.IsAbs(pkgorg) {
			return pkgorg
		}

		if arp := filepath.Join(utils.AppRoot, "vendor", pkgorg); isExistingDir(arp) {
			return arp
		}

		for _, gopath := range utils.GOPATH() {
			if gp := filepath.Join(gopath, "src", pkgorg); isExistingDir(gp) {
				return gp
			}
		}
	}

	return ""
}

func isExistingDir(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}
