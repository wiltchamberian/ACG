package parser

type ParserControllor interface {
	ResetIndex(pr *BasicParser)
	AdvanceIndex(pr *BasicParser)
	Print(pr *Generator)
	PrintNewParser(pr *Generator, name string)
	PrintImportSlice(pr *Generator)

	//farest error position
	UpdateErrorPosition(pos int)
	GetErrorPosition() int

	// AltStart(pr *Generator) int
	// AltEnd(pr *Generator) int
	// GroupStart(pr *Generator) int
	// GroupEnd(pr *Generator) int
}

type ParserControllorL struct {
	errorPos int
}

func (s *ParserControllorL) UpdateErrorPosition(pos int) {
	if s.errorPos < pos {
		s.errorPos = pos
	}
}

func (s *ParserControllorL) GetErrorPosition() int {
	return s.errorPos
}

func (s *ParserControllorL) ResetIndex(pr *BasicParser) {
	pr.index = 0
}

func (s *ParserControllorL) AdvanceIndex(pr *BasicParser) {
	pr.index++
}

func (s *ParserControllorL) Print(pr *Generator) {
}

func (s *ParserControllorL) PrintImportSlice(pr *Generator) {
}

func (s *ParserControllorL) PrintNewParser(pr *Generator, name string) {
	pr.Printf("func New%s() *%s{\n", name, name)
	pr.Printf("\tvar v %s\n", name)
	pr.Printf("\tv.BasicParser = NewBasicParserL()\n")
	pr.Printf("\treturn &v\n")
	pr.Printf("}\n\n")
}

type ParserControllorR struct {
	errorPos int
}

func (s *ParserControllorR) GetErrorPosition() int {
	return s.errorPos
}

func (s *ParserControllorR) UpdateErrorPosition(pos int) {
	if s.errorPos > pos || s.errorPos == 0 {
		s.errorPos = pos
	}
}

func (s *ParserControllorR) ResetIndex(pr *BasicParser) {
	pr.index = len(pr.tokens) - 1
}

func (s *ParserControllorR) AdvanceIndex(pr *BasicParser) {
	pr.index--
}

func (s *ParserControllorR) Print(pr *Generator) {
	pr.Printf("\tslices.Reverse(nodes)\n")
}

func (s *ParserControllorR) PrintNewParser(pr *Generator, name string) {
	pr.Printf("func New%s() *%s{\n", name, name)
	pr.Printf("\tvar v %s\n", name)
	pr.Printf("\tv.BasicParser = NewBasicParserR()\n")
	pr.Printf("\treturn &v\n")
	pr.Printf("}\n\n")
}

func (s *ParserControllorR) PrintImportSlice(pr *Generator) {
	pr.Print("import \"slices\"\n\n")
}
