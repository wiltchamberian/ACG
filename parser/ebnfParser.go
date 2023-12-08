package parser

import "errors"

type Group struct {
	Tokens []*Token
	Type   int //0:item, 1:[]* 2:[]+
}

type Rule struct {
	Name string
	Alts [][]Group
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
			if alt, err := gp.Alternative(); err == nil {
				alts := [][]Group{alt}
				apos := gp.Mark()
				for {
					_, err := gp.ExpectValue("|")
					if err != nil {
						break
					}
					alt, err := gp.Alternative()
					if err != nil {
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
func (gp *EBNFParser) Alternative() ([]Group, error) {
	var groups []Group
	for {
		gp, err := gp.Group()
		if err == nil {
			groups = append(groups, gp)
		} else {
			break
		}
	}
	return groups, nil
}

func (gp *EBNFParser) Group() (group Group, err error) {
	_, err = gp.Expect(TkLBracket)
	if err != nil {
		item := gp.Item()
		if item == nil {
			err = errors.New("Group notmatch")
			return
		} else {
			group.Tokens = append(group.Tokens, item)
			group.Type = 0
			err = nil
			return
		}
	}
	for {
		item := gp.Item()
		if item != nil {
			group.Tokens = append(group.Tokens, item)
		} else {
			break
		}
	}
	_, err = gp.Expect(TkRBracket)
	_, err = gp.Expect(TkMul)
	if err == nil {
		group.Type = 1
		return
	}
	_, err = gp.Expect(TkAdd)
	if err == nil {
		group.Type = 2
		return
	}
	err = errors.New("Group notmatch")
	return
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
