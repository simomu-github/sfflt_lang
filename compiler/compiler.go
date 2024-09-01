package compiler

import (
	"hash/fnv"
	"strings"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

const (
	FUNCTION_LABEL = int64(0b01) << 33

	TMP_ADDR        = int64(0b00)
	GLOBAL_VAR_ADDR = int64(0b01) << 33
	LOCAL_VAR_ADDR  = int64(0b10) << 33
	HEAP_ADDR       = int64(0b11) << 33

	LOCAL_VAR_SCOPE_SHIFT = 8
)

type Compiler struct {
	statements        []ast.Statement
	instructions      instructions
	functions         []instructions
	compilingFunction *compilingFunction
	labelIndex        int
	breakPositions    [][]int
}

type instructions []string

type compilingFunction struct {
	ParamCount int
}

func New(statements []ast.Statement) *Compiler {
	return &Compiler{
		statements:     statements,
		instructions:   instructions{},
		functions:      []instructions{},
		labelIndex:     0,
		breakPositions: [][]int{},
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
	if s.IsLocal {
		addr := intToBinary(LOCAL_VAR_ADDR + int64(s.ScopeDepth<<LOCAL_VAR_SCOPE_SHIFT) + int64(s.LocalIndex))
		c.addInstructionWithParam(PUSH, POSI+addr)
		s.Expression.Visit(c)
		c.addInstruction(STORE)
	} else {
		hash := hashString(s.Identifier.Literal)
		addr := intToBinary(GLOBAL_VAR_ADDR + hash)
		c.addInstructionWithParam(PUSH, POSI+addr)
		s.Expression.Visit(c)
		c.addInstruction(STORE)
	}
}

func (c *Compiler) VisitFunction(s ast.Function) {
	c.compilingFunction = &compilingFunction{ParamCount: len(s.Params)}
	c.functions = append(c.functions, instructions{})

	hash := hashString(s.Name.Literal)
	label := intToBinary(FUNCTION_LABEL + hash)

	c.addInstructionWithParam(LABEL, label)

	for _, stmt := range s.Body {
		stmt.Visit(c)
	}
	c.addInstructionWithParam(PUSH, ZERO)
	if len(s.Params) != 0 {
		slideLength := intToBinary(int64(len(s.Params)))
		c.addInstructionWithParam(SLIDE, POSI+slideLength)
	}
	c.addInstruction(ENDSUB)

	c.compilingFunction = nil
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
	if c.compilingFunction.ParamCount != 0 {
		slideLength := intToBinary(int64(c.compilingFunction.ParamCount))
		c.addInstructionWithParam(SLIDE, POSI+slideLength)
	}
	c.addInstruction(ENDSUB)
}

func (c *Compiler) VisitBreak(s ast.Break) {
	pos := c.reserveJumpLabel(JUMP)
	c.breakPositions[len(c.breakPositions)-1] = append(c.breakPositions[len(c.breakPositions)-1], pos)
}

func (c *Compiler) VisitIf(s ast.If) {
	s.Condition.Visit(c)

	trueJumpPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	s.Then.Visit(c)
	endJumpPos := c.reserveJumpLabel(JUMP)

	trueLabel := c.markJumpLabel()
	c.confirmJumpLabel(trueJumpPos, trueLabel)
	if s.Else != nil {
		s.Else.Visit(c)
	}

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpPos, endLabel)
}

func (c *Compiler) VisitWhile(s ast.While) {
	c.beginLoop()

	trueJumpLabel := c.markJumpLabel()
	s.Condition.Visit(c)
	endJumpPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	s.Body.Visit(c)

	trueJumpPos := c.reserveJumpLabel(JUMP)
	c.confirmJumpLabel(trueJumpPos, trueJumpLabel)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpPos, endLabel)
	breakPositions := c.currentLoopBreakPositions()
	for _, pos := range breakPositions {
		c.confirmJumpLabel(pos, endLabel)
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
	addr := ""
	if s.Target.Type == ast.LOCAL {
		addr = intToBinary(LOCAL_VAR_ADDR + int64(s.Target.ScopeDepth<<LOCAL_VAR_SCOPE_SHIFT) + int64(s.Target.LocalIndex))
	} else {
		hash := hashString(s.Target.Identifier.Literal)
		addr = intToBinary(GLOBAL_VAR_ADDR + hash)
	}
	c.addInstructionWithParam(PUSH, POSI+addr)
	c.addInstruction(RETRIEVE)

	c.addInstruction(DISCARD)

	s.Expression.Visit(c)
	c.addInstruction(DUP)
	c.addInstructionWithParam(PUSH, POSI+addr)
	c.addInstruction(SWAP)
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
		instruction = DIV
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

	zeroJumpPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	if e.Operator.Type == token.EQ {
		c.addInstructionWithParam(PUSH, ZERO)
	} else {
		c.addInstructionWithParam(PUSH, ONE)
	}
	endJumpPos := c.reserveJumpLabel(JUMP)

	zeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(zeroJumpPos, zeroLabel)

	if e.Operator.Type == token.EQ {
		c.addInstructionWithParam(PUSH, ONE)
	} else {
		c.addInstructionWithParam(PUSH, ZERO)
	}

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpPos, endLabel)
}

func (c *Compiler) comparison(e ast.Binary) {
	if e.Operator.Type == token.GT || e.Operator.Type == token.GTEQ {
		c.addInstruction(SWAP)
	}
	c.addInstruction(SUB)

	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		c.addInstruction(DUP)
	}

	zeroJumpPos := -1
	if e.Operator.Type == token.LTEQ || e.Operator.Type == token.GTEQ {
		zeroJumpPos = c.reserveJumpLabel(JUMP_WHEN_ZERO)
	}

	negativeJumpOffset := c.reserveJumpLabel(JUMP_WHEN_NEGA)

	c.addInstructionWithParam(PUSH, ZERO)
	endJumpOffset := c.reserveJumpLabel(JUMP)

	if zeroJumpPos >= 0 {
		zeroLabel := c.markJumpLabel()
		c.confirmJumpLabel(zeroJumpPos, zeroLabel)
		c.addInstruction(DISCARD)
	}

	trueLabel := c.markJumpLabel()
	c.confirmJumpLabel(negativeJumpOffset, trueLabel)
	c.addInstructionWithParam(PUSH, ONE)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpOffset, endLabel)
}

func (c *Compiler) and(e ast.Binary) {
	lhsJumpPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	e.Right.Visit(c)

	rhsJumpPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)

	c.addInstructionWithParam(PUSH, ONE)
	endJumpPos := c.reserveJumpLabel(JUMP)

	zeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(lhsJumpPos, zeroLabel)
	c.confirmJumpLabel(rhsJumpPos, zeroLabel)
	c.addInstructionWithParam(PUSH, ZERO)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(endJumpPos, endLabel)
}

func (c *Compiler) or(e ast.Binary) {
	lhsJumpZeroPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)
	c.addInstructionWithParam(PUSH, ONE)
	lhsJumpEndPos := c.reserveJumpLabel(JUMP)

	lhsJumpZeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(lhsJumpZeroPos, lhsJumpZeroLabel)
	e.Right.Visit(c)

	rhsJumpZeroPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)
	c.addInstructionWithParam(PUSH, ONE)
	rhsJumpEndPos := c.reserveJumpLabel(JUMP)

	rhsJumpZeroLabel := c.markJumpLabel()
	c.confirmJumpLabel(rhsJumpZeroPos, rhsJumpZeroLabel)

	c.addInstructionWithParam(PUSH, ZERO)

	endLabel := c.markJumpLabel()
	c.confirmJumpLabel(lhsJumpEndPos, endLabel)
	c.confirmJumpLabel(rhsJumpEndPos, endLabel)
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
		zeroJumpPos := c.reserveJumpLabel(JUMP_WHEN_ZERO)

		c.addInstructionWithParam(PUSH, ZERO)
		endJumpPos := c.reserveJumpLabel(JUMP)

		zeroLabel := c.markJumpLabel()
		c.confirmJumpLabel(zeroJumpPos, zeroLabel)

		c.addInstructionWithParam(PUSH, ONE)

		endLabel := c.markJumpLabel()
		c.confirmJumpLabel(endJumpPos, endLabel)
		return
	}

	e.Right.Visit(c)
}

func (c *Compiler) VisitCall(e ast.Call) {
	for _, arg := range e.Arguments {
		arg.Visit(c)
	}
	hash := hashString(e.Callee.Literal)
	label := intToBinary(FUNCTION_LABEL + hash)

	c.addInstructionWithParam(CALLSUB, label)
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
	if e.Type == ast.ARGUMENT {
		c.argumentVariable(e)
	} else if e.Type == ast.LOCAL {
		c.localVariable(e)
	} else {
		c.globalVariable(e)
	}
}

func (c *Compiler) argumentVariable(e ast.Variable) {
	offset := c.compilingFunction.ParamCount - e.ArgumentIndex + e.RelativeIndex
	param := intToBinary(int64(offset))
	c.addInstructionWithParam(COPY, POSI+param)
}

func (c *Compiler) globalVariable(e ast.Variable) {
	hash := hashString(e.Identifier.Literal)
	addr := intToBinary(GLOBAL_VAR_ADDR + hash)
	c.addInstructionWithParam(PUSH, POSI+addr)
	c.addInstruction(RETRIEVE)
}

func (c *Compiler) localVariable(e ast.Variable) {
	addr := intToBinary(LOCAL_VAR_ADDR + int64(e.ScopeDepth<<LOCAL_VAR_SCOPE_SHIFT) + int64(e.LocalIndex))
	c.addInstructionWithParam(PUSH, POSI+addr)
	c.addInstruction(RETRIEVE)
}

func (c *Compiler) VisitGet(s ast.Get) {
	addr := intToBinary(TMP_ADDR)
	c.addInstructionWithParam(PUSH, POSI+addr)
	if s.Token.Type == token.GETN {
		c.addInstruction(GETN)
	} else {
		c.addInstruction(GETC)
	}

	c.addInstructionWithParam(PUSH, POSI+addr)
	c.addInstruction(RETRIEVE)
}

func (c *Compiler) VisitArrayLiteral(e ast.ArrayLiteral) {
}

func (c *Compiler) VisitIndex(e ast.Index) {
}

func (c *Compiler) addInstruction(instruction InstructionType) {
	if c.isCompilingFunction() {
		idx := len(c.functions) - 1
		c.functions[idx] = append(c.functions[idx], string(instruction))
	} else {
		c.instructions = append(c.instructions, string(instruction))
	}
}

func (c *Compiler) addInstructionWithParam(instruction InstructionType, param string) {
	if c.isCompilingFunction() {
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
	label := intToBinary(int64(c.labelIndex))
	c.labelIndex++

	c.addInstructionWithParam(LABEL, label)

	return label
}

func (c *Compiler) confirmJumpLabel(pos int, label string) {
	c.currentInstructions()[pos] = strings.Replace(c.currentInstructions()[pos], "?", label, 1)
}

func (c *Compiler) currentInstructions() []string {
	if !c.isCompilingFunction() {
		return c.instructions
	}

	idx := len(c.functions) - 1
	return c.functions[idx]
}

func (c *Compiler) isCompilingFunction() bool {
	return c.compilingFunction != nil
}

func (c *Compiler) beginLoop() {
	c.breakPositions = append(c.breakPositions, []int{})
}

func (c *Compiler) currentLoopBreakPositions() []int {
	return c.breakPositions[len(c.breakPositions)-1]
}

func (c *Compiler) endLoop() {
	c.breakPositions = c.breakPositions[:len(c.breakPositions)-1]
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

func hashString(str string) int64 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return int64(h.Sum32())
}
