package lib

import _ "embed"

//go:embed strings.sflt
var strings string

var buildinLibraries = map[string]string{
	"strings": strings,
}

func LookupBuilinLibrary(name string) (string, bool) {
	lib, ok := buildinLibraries[name]
	return lib, ok
}
