package parser

import "errors"

type Rule struct {
	Name string
	Alts [][]*Token
}

type EBNFParser struct {
	BasicParser
}

// Grammar 方法表示解析整个文法
func (gp *EBNFParser) Grammar() []Rule {
	pos := gp.Mark()
	if rule, err := gp.Rule(); err == nil {
		rules := []Rule{rule}
		for rule, err := gp.Rule(); err == nil; rule, err = gp.Rule() {
			rules = append(rules, rule)
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
func (gp *EBNFParser) Rule() (Rule, error) {
	var rule Rule
	pos := gp.Mark()
	if name, err := gp.Expect(TkIdentifier); err == nil {
		if _, err := gp.Expect(TkColon); err == nil {
			if alt := gp.Alternative(); alt != nil {
				alts := [][]*Token{alt}
				apos := gp.Mark()
				for {
					_, err := gp.ExpectValue("|")
					if err != nil {
						break
					}
					alt := gp.Alternative()
					if len(alt) == 0 {
						break
					}
					alts = append(alts, alt)
					apos = gp.Mark()
				}
				gp.Reset(apos)
				if _, err := gp.Expect(TkSemicolon); err == nil {
					return Rule{string(name.Literal), alts}, err
				}
			}
		}
	}
	gp.Reset(pos)
	return rule, errors.New("no rule")
}

// Alternative 方法表示解析文法规则的一个选择项
func (gp *EBNFParser) Alternative() []*Token {
	items := []*Token{}
	for item := gp.Item(); item != nil; item = gp.Item() {
		items = append(items, item)
	}
	return items
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
