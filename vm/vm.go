package vm

import "fmt"

type Instr uint64

// CPU is a quick and dirty stack-based virtual machine
// which can be used to test code generation without
// having to run an external assembler.
type CPU struct {
	regs       [8]int
	stack      []int
	locals     []int
	localStack [][]int
	callStack  []int

	program []Instr
	pc      int

	Trace bool
}

func NewCPU(prog []Instr) *CPU {
	return &CPU{program: prog}
}

func (c *CPU) Run() int {
	c.stack = make([]int, 0, 100)
	c.pc = 0

	if c.Trace {
		fmt.Println("prog len:", len(c.program))
	}

	for {
		instr := c.program[c.pc]
		c.pc++

		if c.Trace {
			if opcodeHasArg[instr.Opcode()] {
				fmt.Printf("%04d: %s %d\n", c.pc-1, instr.Opcode(), instr.Arg())
			} else {
				fmt.Printf("%04d: %s\n", c.pc-1, instr.Opcode())
			}
		}

		switch instr.Opcode() {
		case Prologue:
			locals := instr.Arg()
			c.locals = make([]int, locals+1)
		case Push:
			c.stack = append(c.stack, c.regs[0])
		case Pop:
			c.regs[instr.Arg()] = c.stack[len(c.stack)-1]
			c.stack = c.stack[:len(c.stack)-1]
		case LoadLocal:
			c.regs[0] = c.locals[instr.Arg()]
		case StoreLocal:
			// this is kind of a hack... the vm can only take one argument
			// this assumes that the arg reg will be stored in the same numbered local
			c.locals[instr.Arg()] = c.regs[instr.Arg()]
		case Load:
			c.regs[0] = c.locals[c.regs[0]]
		case Store:
			c.locals[c.regs[1]] = c.regs[0]
		case LocalAddr:
			c.regs[0] = instr.Arg()
		case LoadInt:
			c.regs[0] = instr.Arg()
		case Add:
			c.regs[0] = c.regs[1] + c.regs[0]
		case Sub:
			c.regs[0] = c.regs[1] - c.regs[0]
		case Mul:
			c.regs[0] = c.regs[1] * c.regs[0]
		case Div:
			c.regs[0] = c.regs[1] / c.regs[0]
		case Neg:
			c.regs[0] = -c.regs[0]
		case Eq:
			if c.regs[1] == c.regs[0] {
				c.regs[0] = 1
			} else {
				c.regs[0] = 0
			}
		case Ne:
			if c.regs[1] != c.regs[0] {
				c.regs[0] = 1
			} else {
				c.regs[0] = 0
			}
		case Lt:
			if c.regs[1] < c.regs[0] {
				c.regs[0] = 1
			} else {
				c.regs[0] = 0
			}
		case Le:
			if c.regs[1] <= c.regs[0] {
				c.regs[0] = 1
			} else {
				c.regs[0] = 0
			}
		case Gt:
			if c.regs[1] > c.regs[0] {
				c.regs[0] = 1
			} else {
				c.regs[0] = 0
			}
		case Ge:
			if c.regs[1] >= c.regs[0] {
				c.regs[0] = 1
			} else {
				c.regs[0] = 0
			}
		case Call:
			c.callStack = append(c.callStack, c.pc)
			c.localStack = append(c.localStack, c.locals)
			c.pc = instr.Arg()
		case JumpIfFalse:
			if c.regs[0] == 0 {
				c.pc = instr.Arg()
			}
		case Jump:
			c.pc = instr.Arg()
		case Return:
			c.pc = c.callStack[len(c.callStack)-1]
			c.callStack = c.callStack[:len(c.callStack)-1]

			c.locals = c.localStack[len(c.localStack)-1]
			c.localStack = c.localStack[:len(c.localStack)-1]
		case Exit:
			return c.regs[0]
		default:
			panic("unknown opcode")
		}
	}
}
