package main

import (
	"fmt"
	"os"
)

const MAX_MEMORY_SIZE = 30_000

const INC = "+"
const DEC = "-"
const PTR_INC = ">"
const PTR_DEC = "<"
const OP_LOOP = "["
const CL_LOOP = "]"
const OUT = "."
const IN = ","

type Parser struct {
	data string
	pos  int
}

func newParser(data string) Parser {
	return Parser{
		data: data,
		pos:  0,
	}
}

func getCurrentParser(p Parser) string {
	return string(p.data[p.pos])
}

func nextPos(p_pos *int) {
	*p_pos += 1
}

type Program struct {
	inst       []string
	inst_pos   int
	memory     []int
	memory_pos int
	loop_idxs  []int
}

func newProg() Program {
	return Program{
		inst_pos:   0,
		memory_pos: 0,
	}
}

func getCurrentInst(prog Program) string {
	return prog.inst[prog.inst_pos]
}

func getCurrentMemory(prog Program) int {
	return prog.memory[prog.memory_pos]
}

func getCurrentLoopIdx(prog Program) int {
	return prog.loop_idxs[len(prog.loop_idxs)-1]
}

func updateMemory(memory *[]int, ptr int) {
	if (ptr+1)-(len(*memory)) > 0 {
		*memory = append(*memory, 0)
	}
}

func updateInstPos(pos *int) {
	*pos += 1
}

func increase(memory []int, ptr int) []int {
	new_mem := memory

	if new_mem[ptr] != 255 {
		new_mem[ptr] += 1
	}

	return new_mem
}

func decrease(memory []int, ptr int) []int {
	new_mem := memory

	if new_mem[ptr] != 0 {
		new_mem[ptr] -= 1
	}

	return new_mem
}

func ptrIncrease(ptr *int) {
	if *ptr == MAX_MEMORY_SIZE {
		*ptr = 0
	} else {
		*ptr += 1
	}
}

func ptrDecrease(ptr *int) {
	if *ptr == 0 {
		*ptr = 255
	} else {
		*ptr -= 1
	}
}

func eof(prog Program) bool {
	return prog.inst_pos >= len(prog.inst)
}

func openLoop(prog Program) Program {
	new_prog := prog
	if getCurrentMemory(new_prog) != 0 {
		new_prog.loop_idxs = append(new_prog.loop_idxs, new_prog.inst_pos)
		return new_prog
	}

	new_prog.loop_idxs = new_prog.loop_idxs[:len(new_prog.loop_idxs)-1]
	num_op_loops := 0
	for getCurrentInst(new_prog) != CL_LOOP && !eof(new_prog) {
		if getCurrentInst(new_prog) == OP_LOOP {
			num_op_loops += 1
		}

		if num_op_loops > 0 && new_prog.inst[new_prog.inst_pos+1] == CL_LOOP {
			new_prog.inst_pos += 1
		}
		new_prog.inst_pos += 1
	}
	return new_prog
}

func closeLoop(prog Program) Program {
	new_prog := prog
	if getCurrentMemory(new_prog) != 0 {
		new_prog.inst_pos = getCurrentLoopIdx(new_prog)
	}
	return new_prog
}

func out(memory_value int) {
	c := rune(memory_value)
	fmt.Printf("%c", c)
}

func main() {
	args := os.Args[1:]
	path := args[0]
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	prog := newProg()
	updateMemory(&prog.memory, prog.memory_pos)

	parser := newParser(string(data))
	for parser.pos < len(string(parser.data)) {
		c := getCurrentParser(parser)
		if checkSyntax(c) {
			prog.inst = append(prog.inst, c)
		}
		nextPos(&parser.pos)
	}

	for eof(prog) {
		inst := getCurrentInst(prog)

		switch inst {
		case INC:
			prog.memory = increase(prog.memory, prog.memory_pos)
			break
		case DEC:
			prog.memory = decrease(prog.memory, prog.memory_pos)
			break
		case PTR_INC:
			ptrIncrease(&prog.memory_pos)
			updateMemory(&prog.memory, prog.memory_pos)
			break
		case PTR_DEC:
			ptrDecrease(&prog.memory_pos)
			updateMemory(&prog.memory, prog.memory_pos)
			break
		case OP_LOOP:
			prog = openLoop(prog)
			break
		case CL_LOOP:
			prog = closeLoop(prog)
			break
		case OUT:
			out(getCurrentMemory(prog))
			break
		case IN:
			break
		default:
			panic("Parsing error")
		}

		prog.inst_pos += 1
	}
}

func checkSyntax(c string) bool {
	return (c == INC || c == DEC ||
		c == PTR_INC || c == PTR_DEC ||
		c == OP_LOOP || c == CL_LOOP ||
		c == OUT || c == IN)
}
