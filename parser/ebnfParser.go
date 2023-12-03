package parser

type Rule struct {
	Name string
	Alts []interface{}
}

type EBNFParser struct {
	BasicParser
}

// Grammar 方法表示解析整个文法
func (gp *EBNFParser) Grammar() []Rule {
	pos := gp.Mark()
	if rule := gp.Rule(); rule != nil {
		rules := []Rule{*rule}
		for rule := gp.Rule(); rule != nil; rule = gp.Rule() {
			rules = append(rules, *rule)
		}
		if _, err := gp.Expect(TkEof); err == nil {
			return rules
		}
		return rules
	}
	gp.Reset(pos)
	return nil
}

// Rule 方法表示解析文法规则
func (gp *EBNFParser) Rule() *Rule {
	pos := gp.Mark()
	if name, err := gp.Expect(TkIdentifier); err == nil {
		if _, err := gp.ExpectValue(":"); err != nil {
			if alt := gp.Alternative(); alt != nil {
				alts := []interface{}{*alt}
				apos := gp.Mark()
				for {
					alt := gp.Alternative()
					_, err := gp.ExpectValue("|")
					if alt == nil || err != nil {
						break
					}

					alts = append(alts, *alt)
					apos = gp.Mark()
				}
				gp.Reset(apos)
				if _, err := gp.Expect(TkSemicolon); err == nil {
					return &Rule{string(name.Literal), alts}
				}
			}
		}
	}
	gp.Reset(pos)
	return nil
}

// Alternative 方法表示解析文法规则的一个选择项
func (gp *EBNFParser) Alternative() *[]interface{} {
	items := []interface{}{}
	for item := gp.Item(); item != nil; item = gp.Item() {
		items = append(items, *item)
	}
	return &items
}

// Item 方法表示解析文法规则的一个项目
func (gp *EBNFParser) Item() *Token {
	if tk, err := gp.Expect(TkIdentifier); err == nil {
		return tk
	}
	if tk, err := gp.Expect(TkString); err == nil {
		return tk
	}
	return nil
}
