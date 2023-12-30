package compiler

import (
	"strings"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

type Compiler struct {
	statements   []ast.Statement
	instructions []string
}

func New(statements []ast.Statement) *Compiler {
	return &Compiler{
		statements:   statements,
		instructions: []string{},
	}
}

func (c *Compiler) Compile() []string {
	for _, e := range c.statements {
		e.Visit(c)
	}

	return c.instructions
}

func (c *Compiler) VisitPutn(s ast.PutnStatement) {
	s.Expression.Visit(c)

	c.instructions = append(c.instructions, "LTFL")
}

func (c *Compiler) VisitExpression(s ast.ExpressionStatement) {
	s.Expression.Visit(c)
	c.instructions = append(c.instructions, "FTT")
}

func (c *Compiler) VisitBinaryExpression(e ast.Binary) {
	e.Left.Visit(c)
	e.Right.Visit(c)
	var instruction string
	switch e.Operator.Type {
	case token.PLUS:
		instruction = "LFFF"
	case token.MINUS:
		instruction = "LFFL"
	case token.ASTERISK:
		instruction = "LFFT"
	case token.SLASH:
		instruction = "LFLF"
	}

	c.instructions = append(c.instructions, instruction)
}

func (c *Compiler) VisitUnaryExpression(e ast.Unary) {
	if e.Operator.Type == token.MINUS {
		c.instructions = append(c.instructions, "FFLLT")
		e.Right.Visit(c)
		c.instructions = append(c.instructions, "LFFT")
		return
	}

	e.Right.Visit(c)
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
