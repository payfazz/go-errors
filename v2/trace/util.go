package trace

import (
	"strings"
)

func inPkg(what, pkg string) bool {
	if !strings.HasPrefix(what, pkg) {
		return false
	}

	if len(what) == len(pkg) {
		return true
	}

	if len(what) > len(pkg) {
		return what[len(pkg)] == '.' || what[len(pkg)] == '/'
	}

	return false
}

// Filter the locations slice by pkg prefix
func FilterByPkgs(locations []Location, pkgs ...string) []Location {
	var ret []Location
	for _, l := range locations {
		if l.InPkg(pkgs...) {
			ret = append(ret, l)
		}
	}
	return ret
}

// return true if this location is in package pkgs
func (l *Location) InPkg(pkgs ...string) bool {
	for _, pkg := range pkgs {
		if inPkg(l.func_, pkg) {
			return true
		}
	}
	return false
}
