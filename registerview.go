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
// 如果 pkgname 和 subpath 都是空字符，则检查 pkgorg 是否是绝对路径,
// 同时会在 $gopath/src, vendor, 等目录下循环检查 pkgorg 项目是否存在。
// 否则会在 $gopath/pkg/mod/ 目录下（go mod 模式）  pkgorg+pkgname+subpath （subpath 为空字符则会默认为 views ,默认会视图目录） 路径是否存在，
// 如果需要检查非 views 目录例如 */*/views 或者 user/theme/local ，则需要设置 subpath 为  */*/views 或者 user/theme/local 。
// 只要检查存在目录就会返回相关路径
func DetectViewsDir(pkgorg, pkgname, subpath string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH

	}

	if pkgname == "" && subpath == "" {
		if abspath := isAbsOrVendorAndSrc(pkgorg); len(abspath) > 0 {
			return abspath
		}
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

// walkPath 循环查询路径
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

// foundPath 返回存在的路径
func foundPath(p, pkgname, subpath string) string {
	if hasPerfix(p, pkgname) {
		if len(subpath) == 0 {
			subpath = "views"
		}
		if vp := filepath.Join(p, subpath); isExistingDir(vp) {
			return vp
		}
	}
	return ""
}

//hasPerfix 判断是否满足前缀
func hasPerfix(p string, pkgname string) bool {
	return strings.HasPrefix(filepath.Base(p), pkgname+"@")
}

//isAbsOrVendorAndSrc 判断路径是否存在于 vendor，src，或者是绝对路径
func isAbsOrVendorAndSrc(pkgorg string) string {

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

// isExistingDir 判断目录是否存在
func isExistingDir(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}
