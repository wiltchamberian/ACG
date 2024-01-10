package parser

const (
	CAT_INTEGER = iota
	CAT_FLOAT
	CAT_STRUCT
)

type FieldShow struct {
	name     string
	typeName string
	offset   int
}

type Field struct {
	name   string
	typ    *NikaType
	offset int
}

type NikaType struct {
	name   string
	fields []FieldShow
	length int

	isReady bool //tmp
	isBasic bool //whether a basic type or not
}

func (s *NikaType) IsInteger() bool {
	return true //TODO
}

type TypeSystem struct {
	typMap map[string]*NikaType
}

func (s *TypeSystem) HasType(name string) bool {
	_, ok := s.typMap[name]
	return ok
}

func (s *TypeSystem) GetType(name string) (*NikaType, bool) {
	result, ok := s.typMap[name]
	return result, ok
}

func (s *TypeSystem) SetType(typ *NikaType) {
	s.typMap[typ.name] = typ
}

func (s *TypeSystem) AddType(tp NikaType) bool {
	if !s.HasType(tp.name) {
		s.typMap[tp.name] = &NikaType{name: tp.name, fields: tp.fields}
		return true
	} else {
		return false
	}
}

func (s *TypeSystem) AlignUp(length int) int {
	return length
}

func (s *TypeSystem) UpdateType(name string) bool {
	tp, ok := s.GetType(name)
	var ok2 bool
	if ok {
		for _, field := range tp.fields {
			ok2 = s.UpdateType(field.name)
			if ok2 == false {
				return false
			}
		}
		offset := 0
		var tmp *NikaType
		for _, field := range tp.fields {
			tmp, ok2 = s.GetType(field.name)
			if ok2 == false {
				return false
			}
			field.typeName = tmp.name
			field.offset = offset
			offset += s.AlignUp(tmp.length)
		}
		tp.isReady = true
		return true
	}
	return false
}

func (s *TypeSystem) UpdateAllTypes() {
	for _, val := range s.typMap {
		if val.isReady == false {
			for _, field := range val.fields {
				s.UpdateType(field.name)
			}
		}
	}
}

func (s *TypeSystem) NewStruct(name string, fields []FieldShow) bool {
	var tp = NikaType{name: name, fields: fields}
	ok := s.AddType(tp)
	return ok
}

func (s *TypeSystem) NewInteger() bool {
	var nk NikaType
	nk.name = "int"
	nk.isBasic = true
	return true
}

func (s *TypeSystem) NewFloat() bool {
	var nk NikaType
	nk.name = "f32"
	nk.isBasic = true
	return true
}
