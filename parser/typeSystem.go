package parser

type Field struct {
	name   string
	typ    *NikaType
	offset int
}

type NikaType struct {
	name   string
	fields []Field
	length int

	isReady bool //tmp
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

func (s *TypeSystem) AddType(tp NikaType) {
	if !s.HasType(tp.name) {
		s.typMap[tp.name] = &NikaType{name: tp.name, fields: tp.fields}
	} else {

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
		for _, field := range tp.fields {
			field.typ, ok2 = s.GetType(field.name)
			if ok2 == false {
				return false
			}
			field.offset = offset
			offset += s.AlignUp(field.typ.length)
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
