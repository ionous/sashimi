package runtime

//
// SourceLookup: so object targets can find their parents
// implemented as a stack to allow context to define hierarchy
//
type SourceLookup func() *GameObject

//
type ParserSourceStack struct {
	arr []SourceLookup
}

//
func (this *ParserSourceStack) FindSource() (ret *GameObject) {
	count := len(this.arr)
	if count > 0 {
		parserSource := this.arr[count-1]
		ret = parserSource()
	}
	return ret
}

//
func (this *ParserSourceStack) PushSource(p SourceLookup) {
	this.arr = append(this.arr, p)
}

//
func (this *ParserSourceStack) PopSource() {
	count := len(this.arr)
	this.arr = this.arr[0 : count-1]
}
