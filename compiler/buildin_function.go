package compiler

var buildinFunctions = map[string]*BuildInFunction{
	"len":    {f: arrayLen, arity: 1},
	"copy":   {f: arrayCopy, arity: 2},
	"append": {f: arrayAppend, arity: 2},

	"_allocate":   {f: allocate, arity: 1},
	"_reallocate": {f: reallocate, arity: 2},
	"_memCopy":    {f: memCopy, arity: 4},
}

type BuildInFunction struct {
	arity int
	f     func(c *Compiler)
}

func arrayLen(c *Compiler) {
	// fetch array pointer
	c.addInstruction(RETRIEVE)
}

// copy(source, dist)
func arrayCopy(c *Compiler) {
	c.addInstructionWithParam(PUSH, ZERO) // counter

	jumpLabel := c.markJumpLabel()

	// arg _memCopy(source, 0, dist, 0)
	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // source
	c.addInstructionWithParam(COPY, POSI+intToBinary(1)) // counter
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // dist
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // counter

	// call _memCopy(source, 0, dist, 0)
	memCopy(c)
	c.addInstruction(DISCARD)

	// update counter
	c.addInstructionWithParam(PUSH, ONE)
	c.addInstruction(ADD)
	c.addInstruction(DUP)
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // source
	c.addInstruction(RETRIEVE)                           // source length
	c.addInstructionWithParam(PUSH, POSI+intToBinary(2))
	c.addInstruction(ADD) // add array header length

	c.addInstruction(SUB) // remaining

	jumpLabelPos := c.reserveJumpLabel(JUMP_WHEN_NEGA)
	c.confirmJumpLabel(jumpLabelPos, jumpLabel)

	// return
	c.addInstructionWithParam(PUSH, ZERO)
	c.addInstructionWithParam(SLIDE, POSI+intToBinary(3))

}

// append(array, element) array
func arrayAppend(c *Compiler) {
	c.addInstructionWithParam(COPY, POSI+intToBinary(1)) // array
	c.addInstructionWithParam(PUSH, ONE)
	c.addInstruction(ADD)
	c.addInstruction(RETRIEVE) // capacity

	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // array
	c.addInstruction(RETRIEVE)                           // length

	c.addInstructionWithParam(PUSH, ONE)
	c.addInstruction(ADD)
	c.addInstruction(SUB) // remain
	reallocJumpPos := c.reserveJumpLabel(JUMP_WHEN_NEGA)

	c.addInstructionWithParam(COPY, POSI+intToBinary(1)) // array
	jumpLabelPos := c.reserveJumpLabel(JUMP)

	// reallocate
	reallocLabel := c.markJumpLabel()
	c.confirmJumpLabel(reallocJumpPos, reallocLabel)

	c.addInstructionWithParam(COPY, POSI+intToBinary(1)) // array

	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // array
	c.addInstructionWithParam(PUSH, ONE)
	c.addInstruction(ADD)
	c.addInstruction(RETRIEVE) // capacity
	c.addInstructionWithParam(PUSH, POSI+intToBinary(2))
	c.addInstruction(MUL) // new capacity

	// call _reallocate(array, capacity)
	reallocate(c)

	jumpLabel := c.markJumpLabel()
	c.confirmJumpLabel(jumpLabelPos, jumpLabel)

	// append
	c.addInstruction(DUP)
	c.addInstruction(RETRIEVE) // length

	c.addInstruction(DUP)                                // length
	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // array
	c.addInstruction(SWAP)
	c.addInstructionWithParam(PUSH, POSI+intToBinary(2))
	c.addInstruction(ADD)
	c.addInstruction(ADD)                                // array last
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // element
	c.addInstruction(STORE)

	c.addInstructionWithParam(PUSH, ONE)
	c.addInstruction(ADD)
	c.addInstructionWithParam(COPY, POSI+intToBinary(1)) // new_array
	c.addInstruction(SWAP)
	c.addInstruction(STORE) // update array length

	// return
	c.addInstructionWithParam(SLIDE, POSI+intToBinary(2))
}

// _memCopy(source, source_index, dist, dist_index)
func memCopy(c *Compiler) {
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // source
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // source_index
	c.addInstruction(ADD)
	c.addInstruction(RETRIEVE)

	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // dist
	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // dist_index
	c.addInstruction(ADD)

	c.addInstruction(SWAP)

	c.addInstruction(STORE)
	c.addInstructionWithParam(PUSH, ZERO)
	c.addInstructionWithParam(SLIDE, POSI+intToBinary(4))
}

// _allocate(size)
func allocate(c *Compiler) {
	c.addInstructionWithParam(PUSH, POSI+intToBinary(VM_ALLOC_REC))
	c.addInstruction(RETRIEVE)
	c.addInstruction(DUP)
	c.addInstructionWithParam(COPY, POSI+intToBinary(2))
	c.addInstruction(ADD)
	c.addInstructionWithParam(PUSH, POSI+intToBinary(VM_ALLOC_REC))
	c.addInstruction(SWAP)
	c.addInstruction(STORE)

	c.addInstructionWithParam(SLIDE, ONE)
}

// _reallocate(original_array, capacity)
func reallocate(c *Compiler) {
	// call _allocate(size)
	c.addInstruction(DUP)
	c.addInstructionWithParam(PUSH, POSI+intToBinary(2))
	c.addInstruction(ADD)
	allocate(c)

	c.addInstruction(DUP)
	c.addInstructionWithParam(COPY, POSI+intToBinary(3)) // original
	c.addInstruction(SWAP)

	// call copy(source, dist)
	arrayCopy(c)
	c.addInstruction(DISCARD)

	// update capacity
	c.addInstruction(DUP)
	c.addInstructionWithParam(PUSH, ONE)
	c.addInstruction(ADD)
	c.addInstructionWithParam(COPY, POSI+intToBinary(2)) // new capacity
	c.addInstruction(STORE)

	// return
	c.addInstructionWithParam(SLIDE, POSI+intToBinary(2))

}
