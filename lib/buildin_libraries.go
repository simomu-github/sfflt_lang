package lib

import _ "embed"

//go:embed strings.sflt
var strings string

//go:embed arrays.sflt
var arrays string

var buildinLibraries = map[string]string{
	"strings": strings,
	"arrays":  arrays,
}

func LookupBuilinLibrary(name string) (string, bool) {
	lib, ok := buildinLibraries[name]
	return lib, ok
}
