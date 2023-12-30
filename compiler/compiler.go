package compiler

import (
	"strings"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

type Compiler struct {
	expressions  []ast.Expression
	instructions []string
}

func New(expressions []ast.Expression) *Compiler {
	return &Compiler{
		expressions:  expressions,
		instructions: []string{},
	}
}

func (c *Compiler) Compile() []string {
	for _, e := range c.expressions {
		e.Visit(c)
	}

	return c.instructions
}

func (c *Compiler) VisitUnaryExpression(e ast.Unary) {
	switch e.Operator.Type {
	case token.MINUS:
		c.instructions = append(c.instructions, "FFLLT")
		e.Right.Visit(c)
		c.instructions = append(c.instructions, "LFFT")
	}
}

func (c *Compiler) VisitIntegerLiteral(e ast.IntegerLiteral) {
	value := intToBinary(e.Value)
	instruction := "FF" + value
	c.instructions = append(c.instructions, instruction)
}

func (c *Compiler) VisitCharLiteral(e ast.CharLiteral) {
	value := intToBinary(int64([]rune(e.Value)[0]))
	instruction := "FF" + value
	c.instructions = append(c.instructions, instruction)
}

func (c *Compiler) VisitBooleanLiteral(e ast.BooleanLiteral) {
	var value string
	if e.Value {
		value = "FLT"
	} else {
		value = "FFT"
	}

	instruction := "FF" + value
	c.instructions = append(c.instructions, instruction)
}

func intToBinary(value int64) string {
	binary := []string{}

	decimal := value
	for decimal != 0 {
		bin := decimal % 2
		if bin == 0 {
			binary = append(binary, "F")
		} else {
			binary = append(binary, "L")
		}
		decimal /= 2
	}

	if value >= 0 {
		binary = append(binary, "F")
	} else {
		binary = append(binary, "L")
	}

	for i := 0; i < len(binary)/2; i++ {
		binary[i], binary[len(binary)-i-1] = binary[len(binary)-i-1], binary[i]
	}
	return strings.Join(binary, "") + "T"
}
