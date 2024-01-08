package parser

type StackMachine struct {
	constants    *DataStack
	instructions Instructions

	*DataStack //variable stack

	//global variable
	globals *DataStack

	//global
	trueVal  int32
	falseVal int32

	//debug
	DebugStack *DataStack
}

func NewStackMachine(compiler *NikaCompiler) *StackMachine {
	var vm StackMachine
	vm.constants = compiler.constants
	vm.instructions = compiler.instructions
	vm.globals = compiler.globals
	vm.DataStack = NewDataStack()
	vm.DebugStack = NewDataStack()

	vm.trueVal = 1
	vm.falseVal = 0
	return &vm
}

func (s *StackMachine) IsEmpty() bool {
	return s.DataStack.Length() <= 0
}

func (s *StackMachine) Print() error {
	file := FileWriter{}
	file.SetOutputPath("./vm.txt")
	file.OpenFile("./vm.txt")
	defer file.CloseFile()
	ip := 0
	l := len(s.instructions)
	for ip < l {
		opCode := s.instructions[ip]
		file.Printf("opcode:%s\n", LookUp(opCode).Name)
		bytelen := InstructionByteLen(opCode)
		ip += bytelen
	}
	file.Flush()
	return nil
}

func (s *StackMachine) Run() error {
	ip := 0
	l := len(s.instructions)
	for ip < l {
		opCode := s.instructions[ip]
		switch opCode {
		case OpConstant:
			{
				index := ReadUint16(s.instructions[ip+1:])
				s.PushInteger(s.constants.GetInteger(int(index)))
			}
		case OpAdd:
			{
				l, r := s.PopTwoIntegers()
				s.PushInteger(l + r)
			}
		case OpSub:
			{
				l, r := s.PopTwoIntegers()
				s.PushInteger(l - r)
			}
		case OpMul:
			{
				l, r := s.PopTwoIntegers()
				s.PushInteger(l * r)
			}
		case OpDiv:
			{
				l, r := s.PopTwoIntegers()
				s.PushInteger(l / r)
			}
		case OpTrue:
			{
				s.PushBool(true)
			}
		case OpFalse:
			{
				s.PushBool(false)
			}
		case OpEq:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool(l == r)
			}
		case OpNotEq:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool(l != r)
			}
		case OpLe:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool(l < r)
			}
		case OpGt:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool(l > r)
			}
		case OpLeEq:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool(l <= r)
			}
		case OpGtEq:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool(l >= r)
			}
		case OpNeg:
			{
				l := s.PopInteger()
				s.PushInteger(-l)
			}
		case OpBang:
			{
				l := s.PopBool()
				s.PushBool(!l)
			}
		case OpOr:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool((l != 0) || (r != 0))
			}
		case OpAnd:
			{
				l, r := s.PopTwoIntegers()
				s.PushBool((l != 0) && (r != 0))
			}
		case OpGlobalSet:
			{
				index := ReadUint16(s.instructions[ip+1:])
				s.globals.SetInteger(int(index), s.TopInteger())
			}
		case OpGlobalGet:
			{
				index := ReadUint16(s.instructions[ip+1:])
				s.PushInteger(s.globals.GetInteger(int(index)))
			}
		case OpJumpNotTrue:
			{
				index := int(ReadUint16(s.instructions[ip+1:]))
				ok := s.PopBool()
				if !ok {
					ip = index - InstructionByteLen(opCode) //cancel the add in the end
				}
			}
		case OpJump:
			{
				ip = int(ReadUint16(s.instructions[ip+1:])) - InstructionByteLen(opCode)
			}
		case OpPop:
			{
				//s.Pop()
				s.PopInteger()
			}
		case OpPopd:
			{
				obj := s.PopInteger()
				s.DebugStack.PushInteger(obj)
			}
		}
		ip += InstructionByteLen(opCode)
	}
	return nil
}
