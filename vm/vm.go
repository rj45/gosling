package vm

type Instr uint64

// CPU is a quick and dirty stack-based virtual machine
// which can be used to test code generation without
// having to run an external assembler.
type CPU struct {
	a      int
	b      int
	stack  []int
	locals []int

	program []Instr
	pc      int
}

func NewCPU(prog []Instr) *CPU {
	return &CPU{program: prog}
}

func (c *CPU) Run() int {
	for {
		instr := c.program[c.pc]
		c.pc++
		switch instr.Opcode() {
		case Prologue:
			locals := instr.Arg()
			c.locals = make([]int, locals+1)
			c.stack = make([]int, 0, 100)
		case Push:
			c.stack = append(c.stack, c.a)
		case Pop:
			c.b = c.stack[len(c.stack)-1]
			c.stack = c.stack[:len(c.stack)-1]
		case LoadLocal:
			c.a = c.locals[instr.Arg()]
		case Load:
			c.a = c.locals[c.a]
		case Store:
			c.locals[c.b] = c.a
		case LocalAddr:
			c.a = instr.Arg()
		case LoadInt:
			c.a = instr.Arg()
		case Add:
			c.a = c.b + c.a
		case Sub:
			c.a = c.b - c.a
		case Mul:
			c.a = c.b * c.a
		case Div:
			c.a = c.b / c.a
		case Neg:
			c.a = -c.a
		case Eq:
			if c.b == c.a {
				c.a = 1
			} else {
				c.a = 0
			}
		case Ne:
			if c.b != c.a {
				c.a = 1
			} else {
				c.a = 0
			}
		case Lt:
			if c.b < c.a {
				c.a = 1
			} else {
				c.a = 0
			}
		case Le:
			if c.b <= c.a {
				c.a = 1
			} else {
				c.a = 0
			}
		case Gt:
			if c.b > c.a {
				c.a = 1
			} else {
				c.a = 0
			}
		case Ge:
			if c.b >= c.a {
				c.a = 1
			} else {
				c.a = 0
			}
		case JumpIfFalse:
			if c.a == 0 {
				c.pc = instr.Arg()
			}
		case Jump:
			c.pc = instr.Arg()
		case Epilogue:
			return c.a
		default:
			panic("unknown opcode")
		}
	}
}
