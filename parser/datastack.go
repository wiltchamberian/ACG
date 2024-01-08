package parser

type DataStack struct {
	stack []byte
}

func NewDataStack() *DataStack {
	var s DataStack
	s.stack = make([]byte, 0)
	return &s
}

func (s *DataStack) PushNode(node INode) {
	if node.GetType().IsInteger() {
		s.PushInteger(node.GetInteger())
	}
	s.PushInteger(node.GetInteger())
}

func (s *DataStack) GetByte(i int) byte {
	return s.stack[i]
}

func (s *DataStack) SetByte(index int, val byte) {
	s.stack[index] = val
}

func (s *DataStack) GetInteger(index int) int32 {
	return ReadInt32(s.stack[index:])
}

func (s *DataStack) SetInteger(index int, val int32) {
	WriteInt32(s.stack[index:], val)
}

func (s *DataStack) GetInt64(index int) int64 {
	return ReadInt64(s.stack[index:])
}

func (s *DataStack) SetInt64(index int, val int64) {
	WriteInt64(s.stack[index:], val)
}

func (s *DataStack) GetUint64(index int) uint64 {
	return ReadUint64(s.stack[index:])
}

func (s *DataStack) SetUint64(index int, val uint64) {
	WriteUint64(s.stack[index:], val)
}

func (s *DataStack) GetUint16(index int) uint16 {
	return ReadUint16(s.stack[index:])
}

func (s *DataStack) SetUint16(index int, val uint16) {
	WriteUint16(s.stack[index:], val)
}

func (s *DataStack) GetInt16(index int) int16 {
	return ReadInt16(s.stack[index:])
}

func (s *DataStack) SetInt16(index int, val int16) {
	WriteInt16(s.stack[index:], val)
}

func (s *DataStack) SetFloat32(index int, val float32) {
	WriteFloat32(s.stack[index:], val)
}

func (s *DataStack) SetFloat64(index int, val float64) {
	WriteFloat64(s.stack[index:], val)
}

func (s *DataStack) Length() int {
	return len(s.stack)
}

func (s *DataStack) PopTwoIntegers() (int32, int32) {
	right := s.PopInteger()
	left := s.PopInteger()
	return left, right
}

func (s *DataStack) PopInteger() int32 {
	l := len(s.stack)
	x := ReadInt32(s.stack[l-4:])
	s.stack = s.stack[0 : l-4]
	return x
}

func (s *DataStack) PushInteger(val int32) int {
	s.stack = append(s.stack, 0)
	s.stack = append(s.stack, 0)
	s.stack = append(s.stack, 0)
	s.stack = append(s.stack, 0)
	WriteInt32(s.stack[len(s.stack)-4:], val)
	return len(s.stack) - 4
}

func (s *DataStack) PushBool(val bool) {
	if val {
		s.PushInteger(1)
	} else {
		s.PushInteger(0)
	}
}

func (s *DataStack) PopBool() bool {
	t := s.PopInteger()
	if t == 0 {
		return false
	} else {
		return true
	}
}

func (s *DataStack) PushFloat32(val float32) {
	s.PushInteger(int32(val))
}

func (s *DataStack) PopFloat32() float32 {
	return float32(s.PopInteger())
}

func (s *DataStack) TopInteger() int32 {
	return ReadInt32(s.stack[len(s.stack)-4:])
}

func (s *DataStack) TopFloat32() float32 {
	return ReadFloat32(s.stack[len(s.stack)-4:])
}

func (s *DataStack) TopInt64() int64 {
	return ReadInt64(s.stack[len(s.stack)-8:])
}

func (s *DataStack) TopFloat64() float64 {
	return ReadFloat64(s.stack[len(s.stack)-8:])
}
