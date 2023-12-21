package parser

import (
	"github.com/llir/llvm/ir"
)

type IRGenerator struct {
	module *ir.Module
}

func NewIRGenerator() IRGenerator {
	var irGenerator IRGenerator
	irGenerator.module = ir.NewModule()
	return irGenerator
}

func (s *IRGenerator) Terminal(inode INode) {

}

func (s *IRGenerator) Nonterminal(inode INode) {
	name := inode.GetName()

	switch name {
	case "struct":
		{
			// children := inode.GetChildren()
			// structName := children[1].GetName()
			// members := children[3].GetChildren()
			// var lvStruct *types.Struct
			// fields := []*types.Var{}

			// tt := s.module.NewTypeDef()
			// var ttt types.Type

			// for member := range members {
			// 	lis := member.GetChildren()
			// 	fieldName := lis[0]
			// 	typeName := lis[1]
			// }
		}
	}
}
