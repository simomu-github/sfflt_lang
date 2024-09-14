package compiler

var buildinFunctions = map[string]*BuildInFunction{
	"len": {f: arrayLen, arity: 1},
}

type BuildInFunction struct {
	arity int
	f     func(c *Compiler)
}

func arrayLen(c *Compiler) {
	// fetch array pointer
	c.addInstruction(RETRIEVE)
}
