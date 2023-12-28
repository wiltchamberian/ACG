package parser

import (
	"errors"
	"fmt"
)

const (
	GroupType_Item = iota
	GroupType_Kleen
	GroupType_Add
	GroupType_Action
)

type Group struct {
	Tokens []*Token
	Type   int //0:item, 1:[]* 2:[]+
}

type Alter struct {
	action           string
	NotDefaultAction bool
	Groups           []Group
}

type Rule struct {
	Name string
	Alts []Alter
}

type BNFConfig struct {
	DefaultAction string
}

type BNFRules struct {
	Rules       []Rule
	Config      BNFConfig
	ActionAlias map[string]string
}

type BNFParser struct {
	BasicParser
	Config BNFConfig

	ActionAlias map[string]string
}

func NewBNFParser() BNFParser {
	var ebnf = BNFParser{BasicParser: NewBasicParserL()}
	ebnf.BasicParser.lexer.mode = Lexer_Grammar
	ebnf.ActionAlias = make(map[string]string)
	return ebnf
}

func (s *BNFParser) updateActions(rules []Rule) {
	for _, rule := range rules {
		for _, alt := range rule.Alts {
			alt.action = s.updateAction(alt.action)
		}
	}
}

func (s *BNFParser) updateAction(action string) string {
	var output string
	l := len(action)
	for i := 0; i < len(action); i++ {
		if action[i] == '#' {
			start := i
			end := start + 1
			for end < l {
				if action[end] == '#' {
					var key = action[start+1 : end]
					output = output + s.ActionAlias[key]
					i = end
					break
				}
				end++
			}
		} else {
			output = output + string(action[i])
		}
	}
	return output
}

func (gp *BNFParser) Parse() ([]Rule, error) {
	gp.ParseTerminators()
	gp.ParseBNFConfig()
	gp.ParseBNFActionAlias()
	return gp.Grammar()
}

func (gp *BNFParser) ParseTerminators() error {
	_, e := gp.Expect(TkLBracket)
	if e != nil {
		return e
	}
	item, e := gp.Expect(TkIdentifier)
	if e != nil {
		return e
	}
	if string(item.Literal) != "terminations" {
		return errors.New("no terminations")
	}
	_, e = gp.Expect(TkRBracket)
	if e != nil {
		return e
	}

	//
	return nil
}

func (gp *BNFParser) ParseBNFConfig() error {
	_, e := gp.Expect(TkLBracket)
	if e != nil {
		return e
	}
	item, e := gp.Expect(TkIdentifier)
	if e != nil {
		return e
	}
	if string(item.Literal) != "bnf_config" {
		return errors.New("no bnf_config")
	}
	_, e = gp.Expect(TkRBracket)
	if e != nil {
		return e
	}

	item, e = gp.Expect(TkIdentifier)
	if e != nil {
		return e
	}
	if string(item.Literal) != "default_action" {
		return errors.New("no default action")
	}
	_, e = gp.Expect(TkColon)
	if e != nil {
		return e
	}
	item, e = gp.Expect(TkAction)
	if e != nil {
		return e
	}
	gp.Config.DefaultAction = item.GetLiteral()

	return nil
}

func (gp *BNFParser) ParseBNFActionAlias() error {
	_, e := gp.Expect(TkLBracket)
	if e != nil {
		return e
	}
	item, e := gp.Expect(TkIdentifier)
	if e != nil {
		return e
	}
	if string(item.Literal) != "bnf_action_alias" {
		return errors.New("no bnf_action_alias")
	}
	_, e = gp.Expect(TkRBracket)
	if e != nil {
		return e
	}

	for {
		item, e = gp.Expect(TkString)
		if e != nil {
			return e
		}
		_, e = gp.ExpectValue(":")
		if e != nil {
			return e
		}
		content, e := gp.Expect(TkString)
		if e != nil {
			return e
		}
		key := item.GetLiteral()
		key = key[1 : len(key)-1]
		value := content.GetLiteral()
		value = value[1 : len(value)-1]
		gp.ActionAlias[key] = value

	}

}

// Grammar 方法表示解析整个文法
func (gp *BNFParser) Grammar() ([]Rule, error) {
	_, er := gp.Expect(TkLBracket)
	if er != nil {
		return nil, er
	}
	_, er = gp.ExpectValue("bnf")
	if er != nil {
		return nil, er
	}
	_, er = gp.Expect(TkRBracket)
	if er != nil {
		return nil, er
	}

	pos := gp.Mark()
	if rule, err := gp.Rule(); err == nil {
		rules := []Rule{rule}
		for rule, err := gp.Rule(); err == nil; rule, err = gp.Rule() {
			rules = append(rules, rule)
		}
		if _, err := gp.Expect(TkEof); err == nil {
			return rules, err
		}
		return rules, err
	}
	gp.Reset(pos)
	return nil, nil
}

// Rule 方法表示解析文法规则
func (gp *BNFParser) Rule() (Rule, error) {
	var rule Rule
	pos := gp.Mark()
	if name, err := gp.Expect(TkIdentifier); err == nil {
		if _, err := gp.Expect(TkColon); err == nil {
			if alt, err := gp.Alternative(); err == nil {
				alts := []Alter{alt}
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
func (gp *BNFParser) Alternative() (Alter, error) {
	var alt Alter
	alt.action = gp.Config.DefaultAction
	for {
		gp, err := gp.Group()
		if err == nil {
			if gp.Type == GroupType_Action {
				alt.action = string(gp.Tokens[0].Literal)
				alt.NotDefaultAction = true
			} else {
				alt.Groups = append(alt.Groups, gp)
			}
		} else {
			break
		}
	}
	return alt, nil
}

func (gp *BNFParser) Group() (group Group, err error) {
	_, err = gp.Expect(TkLParen)
	if err == nil {
		for {
			item := gp.ItemInGroup()
			if item != nil {
				fmt.Println(string(item.Literal))
				group.Tokens = append(group.Tokens, item)
			} else {
				break
			}
		}
		_, err = gp.Expect(TkRParen)
		_, err = gp.Expect(TkMul)
		if err == nil {
			group.Type = GroupType_Kleen
			return
		}
		_, err = gp.Expect(TkAdd)
		if err == nil {
			group.Type = GroupType_Add
			return
		}
	}
	action, err := gp.Expect(TkAction)
	if err == nil {
		group.Type = GroupType_Action
		group.Tokens = append(group.Tokens, action)
		return
	}

	item := gp.Item()
	if item == nil {
		err = errors.New("Group notmatch")
		return
	} else {
		group.Tokens = append(group.Tokens, item)
		group.Type = GroupType_Item
		err = nil
		return
	}

	return
}

// Item 方法表示解析文法规则的一个项目
func (gp *BNFParser) Item() *Token {
	if tk, err := gp.Expect(TkIdentifier); err == nil {
		return tk
	}
	if tk, err := gp.Expect(TkString); err == nil {
		return tk
	}
	return nil
}

func (gp *BNFParser) ItemInGroup() *Token {
	if tk, err := gp.Expect(TkIdentifier); err == nil {
		return tk
	}
	if tk, err := gp.Expect(TkString); err == nil {
		return tk
	}
	if tk, err := gp.ExpectValue("|"); err == nil {
		return tk
	}
	return nil
}

func (gp *BNFParser) ParseFile(path string) (BNFRules, error) {
	gp.ReadFile(path)
	gp.TokenStream()

	//print tokens to file
	var writer FileWriter
	writer.OpenFile("./tokens.txt")
	for _, tk := range gp.tokens {
		writer.FPrintf("%s\n", string(tk.Literal))
	}
	writer.Flush()
	writer.CloseFile()

	rules, err := gp.Parse()
	gp.updateActions(rules)
	var bnfRules BNFRules
	bnfRules.Rules = rules
	bnfRules.Config = gp.Config
	bnfRules.ActionAlias = gp.ActionAlias
	return bnfRules, err
}
