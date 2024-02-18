package compiler

import (
	"strings"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

type Compiler struct {
	statements          []ast.Statement
	instructions        []string
	functions           [][]string
	labelIndex          int
	isCompilingFunction bool
	breakOffsets        [][]int
}

func New(statements []ast.Statement) *Compiler {
	return &Compiler{
		statements:          statements,
		instructions:        []string{},
		functions:           [][]string{},
		labelIndex:          0,
		isCompilingFunction: false,
		breakOffsets:        [][]int{},
	}
}

func (c *Compiler) Compile() []string {
	for _, e := range c.statements {
		e.Visit(c)
	}

	c.addInstruction(END)

	for _, function := range c.functions {
		for _, inst := range function {
			c.instructions = append(c.instructions, inst)
		}
	}

	return c.instructions
}

func (c *Compiler) VisitVar(s ast.Var) {
	ident := stringToBinary("gv" + s.Identifier.Literal)
	c.addInstructionWithParam(PUSH, POSI+ident)
	s.Expression.Visit(c)
	c.addInstruction(STORE)
}

func (c *Compiler) VisitFunction(s ast.Function) {
	c.isCompilingFunction = true

	c.functions = append(c.functions, []string{})
	ident := stringToBinary("gf" + s.Name.Literal)
	c.addInstructionWithParam(LABEL, ident)

	for _, stmt := range s.Body {
		stmt.Visit(c)
	}
	c.addInstructionWithParam(PUSH, ZERO)
	c.addInstruction(ENDSUB)

	c.isCompilingFunction = false
}

func (c *Compiler) VisitPut(s ast.PutStatement) {
	s.Expression.Visit(c)

	if s.Token.Type == token.PUTN {
		c.addInstruction(PUTN)
	} else {
		c.addInstruction(PUTC)
	}
}

func (c *Compiler) VisitReturn(s ast.Return) {
	s.Value.Visit(c)
	c.addInstruction(ENDSUB)
}

func (c *Compiler) VisitBreak(s ast.Break) {
	c.addBreak()
}

func (c *Compiler) VisitIf(s ast.If) {
	s.Condition.Visit(c)

	trueJumpOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	s.Then.Visit(c)
	endJumpOffset := c.reserveJumpLabel(JUMP)

	trueLabel := c.markJumpLabel()
	c.confirmJumpLabel(trueJumpOffset, trueLabel)
	if s.Else != nil {
		s.Else.Visit(c)
	}

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpOffset, endLabel)
}

func (c *Compiler) VisitWhile(s ast.While) {
	c.beginLoop()

	trueJumpLabel := c.markJumpLabel()
	s.Condition.Visit(c)
	endJumpOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	s.Body.Visit(c)

	trueJumpOffset := c.reserveJumpLabel(JUMP)
	c.confirmJumpLabel(trueJumpOffset, trueJumpLabel)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpOffset, endLabel)
	breakOffsets := c.currentBreakOffsets()
	for _, offset := range breakOffsets {
		c.confirmJumpLabel(offset, endLabel)
	}

	c.endLoop()
}

func (c *Compiler) VisitBlock(s ast.Block) {
	for _, stmt := range s.Statements {
		stmt.Visit(c)
	}
}

func (c *Compiler) VisitExpression(s ast.ExpressionStatement) {
	s.Expression.Visit(c)
	c.addInstruction(DISCARD)
}

func (c *Compiler) VisitAssign(s ast.Assign) {
	ident := stringToBinary("gv" + s.Target.Literal)
	c.addInstructionWithParam(PUSH, POSI+ident)
	c.addInstruction(RETRIEVE)

	c.addInstruction(DISCARD)

	c.addInstructionWithParam(PUSH, POSI+ident)
	s.Expression.Visit(c)
	c.addInstruction(STORE)
}

func (c *Compiler) VisitBinaryExpression(e ast.Binary) {
	e.Left.Visit(c)
	var instruction InstructionType
	switch e.Operator.Type {
	case token.PLUS:
		instruction = ADD
	case token.MINUS:
		instruction = SUB
	case token.ASTERISK:
		instruction = MUL
	case token.SLASH:
		instruction = SUB
	case token.MOD:
		instruction = MOD
	case token.LT, token.LTEQ, token.GT, token.GTEQ:
		e.Right.Visit(c)
		c.comparison(e)
		return
	case token.EQ, token.NOT_EQ:
		e.Right.Visit(c)
		c.equality(e)
		return
	case token.AND:
		c.and(e)
		return
	case token.OR:
		c.or(e)
		return
	}

	e.Right.Visit(c)
	c.addInstruction(instruction)
}

func (c *Compiler) equality(e ast.Binary) {
	c.addInstruction(SUB)

	zeroJumpOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	if e.Operator.Type == token.EQ {
		c.addInstructionWithParam(PUSH, ZERO)
	} else {
		c.addInstructionWithParam(PUSH, ONE)
	}
	endJumpOffset := c.reserveJumpLabel(JUMP)

	zeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(zeroJumpOffset, zeroLabel)

	if e.Operator.Type == token.EQ {
		c.addInstructionWithParam(PUSH, ONE)
	} else {
		c.addInstructionWithParam(PUSH, ZERO)
	}

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpOffset, endLabel)
}

func (c *Compiler) comparison(e ast.Binary) {
	if e.Operator.Type == token.GT || e.Operator.Type == token.GTEQ {
		c.addInstruction(SWAP)
	}
	c.addInstruction(SUB)

	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		c.addInstruction(DUP)
	}

	zeroJumpOffset := -1
	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		zeroJumpOffset = c.reserveJumpLabel(JUMP_WHEN_ZERO)
	}

	negativeJumpOffset := c.reserveJumpLabel(JUMP_WHEN_NEGA)

	c.addInstructionWithParam(PUSH, ZERO)
	endJumpOffset := c.reserveJumpLabel(JUMP)

	if zeroJumpOffset >= 0 {
		zeroLabel := c.markJumpLabel()
		c.confirmJumpLabel(zeroJumpOffset, zeroLabel)
		c.addInstruction(DISCARD)
	}

	trueLabel := c.markJumpLabel()
	c.confirmJumpLabel(negativeJumpOffset, trueLabel)
	c.addInstructionWithParam(PUSH, ONE)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpOffset, endLabel)
}

func (c *Compiler) and(e ast.Binary) {
	lhsJumpOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	e.Right.Visit(c)

	rhsJumpOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	c.addInstructionWithParam(PUSH, ONE)
	endJumpOffset := c.reserveJumpLabel(JUMP)

	zeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(lhsJumpOffset, zeroLabel)
	c.confirmJumpLabel(rhsJumpOffset, zeroLabel)
	c.addInstructionWithParam(PUSH, ZERO)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpOffset, endLabel)
}

func (c *Compiler) or(e ast.Binary) {
	lhsJumpZeroOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)
	c.addInstructionWithParam(PUSH, ONE)
	lhsJumpEndOffset := c.reserveJumpLabel(JUMP)

	lhsJumpZeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(lhsJumpZeroOffset, lhsJumpZeroLabel)
	e.Right.Visit(c)

	rhsJumpZeroOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)
	c.addInstructionWithParam(PUSH, ONE)
	rhsJumpEndOffset := c.reserveJumpLabel(JUMP)

	rhsJumpZeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(rhsJumpZeroOffset, rhsJumpZeroLabel)

	c.addInstructionWithParam(PUSH, ZERO)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(lhsJumpEndOffset, endLabel)
	c.confirmJumpLabel(rhsJumpEndOffset, endLabel)
}

func (c *Compiler) VisitUnaryExpression(e ast.Unary) {
	if e.Operator.Type == token.MINUS {
		c.addInstructionWithParam(PUSH, MINUS_ONE)
		e.Right.Visit(c)
		c.addInstruction(MUL)
		return
	}

	if e.Operator.Type == token.BANG {
		e.Right.Visit(c)
		zeroJumpOffset := c.reserveJumpLabel(JUMP_WHEN_ZERO)

		c.addInstructionWithParam(PUSH, ZERO)
		endJumpOffset := c.reserveJumpLabel(JUMP)

		zeroLabel := c.markJumpLabel()
		c.confirmJumpLabel(zeroJumpOffset, zeroLabel)

		c.addInstructionWithParam(PUSH, ONE)

		endLabel := c.markJumpLabel()
		c.confirmJumpLabel(endJumpOffset, endLabel)
		return
	}

	e.Right.Visit(c)
}

func (c *Compiler) VisitCall(e ast.Call) {
	ident := stringToBinary("gf" + e.Callee.Literal)
	c.addInstructionWithParam(CALLSUB, ident)
}

func (c *Compiler) VisitIntegerLiteral(e ast.IntegerLiteral) {
	value := intToBinary(e.Value)
	var sign string
	if e.Value >= 0 {
		sign = POSI
	} else {
		sign = NEGA
	}
	c.addInstructionWithParam(PUSH, sign+value)
}

func (c *Compiler) VisitCharLiteral(e ast.CharLiteral) {
	value := intToBinary(int64([]rune(e.Value)[0]))
	c.addInstructionWithParam(PUSH, POSI+value)
}

func (c *Compiler) VisitBooleanLiteral(e ast.BooleanLiteral) {
	if e.Value {
		c.addInstructionWithParam(PUSH, ONE)
	} else {
		c.addInstructionWithParam(PUSH, ZERO)
	}
}

func (c *Compiler) VisitVariable(e ast.Variable) {
	ident := stringToBinary("gv" + e.Identifier.Literal)
	c.addInstructionWithParam(PUSH, POSI+ident)
	c.addInstruction(RETRIEVE)
}

func (c *Compiler) VisitGet(s ast.Get) {
	tmp := stringToBinary("t" + "tmp")
	c.addInstructionWithParam(PUSH, POSI+tmp)
	if s.Token.Type == token.GETN {
		c.addInstruction(GETN)
	} else {
		c.addInstruction(GETC)
	}

	c.addInstructionWithParam(PUSH, POSI+tmp)
	c.addInstruction(RETRIEVE)
}

func (c *Compiler) addInstruction(instruction InstructionType) {
	if c.isCompilingFunction {
		idx := len(c.functions) - 1
		c.functions[idx] = append(c.functions[idx], string(instruction))
	} else {
		c.instructions = append(c.instructions, string(instruction))
	}
}

func (c *Compiler) addInstructionWithParam(instruction InstructionType, param string) {
	if c.isCompilingFunction {
		idx := len(c.functions) - 1
		c.functions[idx] = append(c.functions[idx], string(instruction)+param+"T")
	} else {
		c.instructions = append(c.instructions, string(instruction)+param+"T")
	}
}

func (c *Compiler) reserveJumpLabel(instruction InstructionType) int {
	c.addInstructionWithParam(instruction, "?")
	return len(c.currentInstructions()) - 1
}

func (c *Compiler) markJumpLabel() string {
	labelPrefix := stringToBinary("l")
	label := intToBinary(int64(c.labelIndex))
	c.labelIndex++

	c.addInstructionWithParam(LABEL, labelPrefix+label)

	return labelPrefix + label
}

func (c *Compiler) confirmJumpLabel(offset int, label string) {
	c.currentInstructions()[offset] = strings.Replace(c.currentInstructions()[offset], "?", label, 1)
}

func (c *Compiler) currentInstructions() []string {
	if !c.isCompilingFunction {
		return c.instructions
	}

	idx := len(c.functions) - 1
	return c.functions[idx]
}

func (c *Compiler) beginLoop() {
	c.breakOffsets = append(c.breakOffsets, []int{})
}

func (c *Compiler) addBreak() {
	offset := c.reserveJumpLabel(JUMP)
	c.breakOffsets[len(c.breakOffsets)-1] = append(c.breakOffsets[len(c.breakOffsets)-1], offset)
}

func (c *Compiler) currentBreakOffsets() []int {
	return c.breakOffsets[len(c.breakOffsets)-1]
}

func (c *Compiler) endLoop() {
	c.breakOffsets = c.breakOffsets[:len(c.breakOffsets)-1]
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
	for {
		bin := decimal % 2
		if bin == 0 {
			binary = append(binary, "F")
		} else {
			binary = append(binary, "L")
		}
		decimal /= 2

		if decimal == 0 {
			break
		}
	}

	for i := 0; i < len(binary)/2; i++ {
		binary[i], binary[len(binary)-i-1] = binary[len(binary)-i-1], binary[i]
	}
	return strings.Join(binary, "")
}
