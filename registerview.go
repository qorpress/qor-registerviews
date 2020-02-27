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

		return ""
	}

	path := filepath.Join(gopath, "/pkg/mod/")
	ppath := filepath.Join(path, pkgorg)
	if _, err := os.Stat(ppath); err == nil {
		var foundp string
		var found bool
		if err = filepath.Walk(ppath, func(p string, f os.FileInfo, err error) error { // nolint: errcheck, gosec, unparam
			if found {
				return nil
			}

			//"media", /media_library,asset_manager
			//"l10n",  /publish
			//"i18n",  /exchange_actions, inline_edit
			//"auth", /auth/providers/password ,facebook ,twitter
			if pkgname != "" {
				if strings.HasPrefix(filepath.Base(p), pkgname+"@") {
					if subpath == "" {
						if vp := filepath.Join(p, "views"); isExistingDir(vp) {
							foundp = vp
							found = true
						} else {
							return err
						}
					} else if subpath != "" {
						if strings.HasPrefix(filepath.Base(p), pkgname+"@") {
							if cvp := filepath.Join(p, subpath+"views"); isExistingDir(cvp) { // auth@
								foundp = cvp
								found = true
							} else {
								return err
							}
						}
					}
				}
			}

			return nil
		}); err != nil {
			color.Red(fmt.Sprintf("os.Stat2 error %v\n", err))
		}

		return foundp
	} else {
		color.Red(fmt.Sprintf("os.Stat1 error %v\n", err))
	}

	return ""
}

func isExistingDir(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}
