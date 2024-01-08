//use an array to represent the
//operation of left-deriviation of multiTree

package parser

type Item struct {
	rule     int
	expanded bool
}

type MultiTree struct {
	arr     []Item
	answers []INode
}

func (s *MultiTree) CreateNode(rule int) INode {
	return nil
}

func (s *MultiTree) GetSubRules(rule int) []Item {
	return []Item{}
}

func (s *MultiTree) GetSubRuleCount(rule int) int {
	return 0
}

func (s *MultiTree) IsTermination(tag int) bool {
	return true
}

func (s *MultiTree) LeftMostDerivate() {
label:
	index := len(s.arr) - 1
	if index < 0 {
		return
	}
	back := s.arr[index].rule
	if s.IsTermination(back) {
		s.answers = append(s.answers, s.CreateNode(back))
		s.arr = s.arr[0:index]
	} else {
		if s.arr[index].expanded == false {
			s.arr[index].expanded = true
			s.arr = append(s.arr, s.GetSubRules(back)...)
			goto label
		} else {
			count := s.GetSubRuleCount(back)
			node := s.CreateNode(back)
			l := len(s.answers)
			for i := len(s.answers) - count; i < len(s.answers); i++ {
				s.answers[i].SetParent(node)
				node.AddChild(s.answers[i])
			}
			s.answers = s.answers[0 : l-count]
			s.answers = append(s.answers, node)
			s.arr = s.arr[0:index]
		}
	}
}
