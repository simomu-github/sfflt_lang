package type_checker

import "github.com/simomu-github/sfflt_lang/token"

type Type struct {
	Tag string
}

type DeclaredType struct {
	Name token.Token
}
