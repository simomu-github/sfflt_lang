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

func (c *Compiler) VisitVar(s ast.Var) {
	ident := stringToBinary("g" + s.Identifier.Literal)
	c.instructions = append(c.instructions, "FFF"+ident+"T")
	s.Expression.Visit(c)
	c.instructions = append(c.instructions, "LLF")
}

func (c *Compiler) VisitPut(s ast.PutStatement) {
	s.Expression.Visit(c)

	if s.Token.Type == token.PUTN {
		c.instructions = append(c.instructions, "LTFL")
	} else {
		c.instructions = append(c.instructions, "LTFF")
	}
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
	case token.LT, token.LTEQ, token.GT, token.GTEQ:
		c.comparison(e)
		return
	}

	c.instructions = append(c.instructions, instruction)
}

func (c *Compiler) comparison(e ast.Binary) {
	if e.Operator.Type == token.GT || e.Operator.Type == token.GTEQ {
		c.instructions = append(c.instructions, "FTL")
	}
	c.instructions = append(c.instructions, "LFFL")

	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		c.instructions = append(c.instructions, "FTF")
	}

	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		// c.instructions = append(c.instructions, "TLF"+labelPrefix+trueLabel+"T")
	}

	whenNegative := c.reserveJumpLabel("TLL")

	c.instructions = append(c.instructions, "FFFFT")
	jumpOffset := c.reserveJumpLabel("TFT")

	trueLabel := c.markJumpLabel()
	c.confirmJumpLabel(whenNegative, trueLabel)
	c.instructions = append(c.instructions, "FFFLT")

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(jumpOffset, endLabel)

	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		c.instructions = append(c.instructions, "FTT")
	}
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
	var sign string
	if e.Value >= 0 {
		sign = "F"
	} else {
		sign = "L"
	}
	instruction := "FF" + sign + value + "T"
	c.instructions = append(c.instructions, instruction)
}

func (c *Compiler) VisitCharLiteral(e ast.CharLiteral) {
	value := intToBinary(int64([]rune(e.Value)[0]))
	instruction := "FFF" + value + "T"
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

func (c *Compiler) VisitVariable(e ast.Variable) {
	ident := stringToBinary("g" + e.Identifier.Literal)
	c.instructions = append(c.instructions, "FFF"+ident+"T")
	c.instructions = append(c.instructions, "LLL")
}

func (c *Compiler) reserveJumpLabel(instruction string) int {
	c.instructions = append(c.instructions, instruction+"?T")
	return len(c.instructions) - 1
}

func (c *Compiler) markJumpLabel() string {
	labelPrefix := stringToBinary("cl")
	label := intToBinary(int64(len(c.instructions)))
	c.instructions = append(c.instructions, "TFF"+labelPrefix+label+"T")

	return labelPrefix + label
}

func (c *Compiler) confirmJumpLabel(offset int, label string) {
	c.instructions[offset] = strings.Replace(c.instructions[offset], "?", label, 1)
}

func stringToBinary(ident string) string {
	result := ""
	for _, char := range ident {
		result += intToBinary(int64(char))
	}

	return result
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

	for i := 0; i < len(binary)/2; i++ {
		binary[i], binary[len(binary)-i-1] = binary[len(binary)-i-1], binary[i]
	}
	return strings.Join(binary, "")
}
