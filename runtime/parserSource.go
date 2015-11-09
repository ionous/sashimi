package runtime

import "github.com/ionous/sashimi/runtime/api"

// SourceLookup: so object targets can find their parents
// implemented as a stack to allow context to define hierarchy
type SourceLookup func() api.Instance

//
type ParserSourceStack struct {
	arr []SourceLookup
}

//
func (s *ParserSourceStack) FindSource() (ret api.Instance) {
	count := len(s.arr)
	if count > 0 {
		parserSource := s.arr[count-1]
		ret = parserSource()
	}
	return ret
}

//
func (s *ParserSourceStack) PushSource(p SourceLookup) {
	s.arr = append(s.arr, p)
}

//
func (s *ParserSourceStack) PopSource() {
	count := len(s.arr)
	s.arr = s.arr[0 : count-1]
}
