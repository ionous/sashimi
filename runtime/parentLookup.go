package runtime

import "github.com/ionous/sashimi/runtime/api"

// TargetLookup: so object targets can find their parents
// implemented as a stack to allow context to define hierarchy
type TargetLookup func(gobj api.Instance) api.Instance

//
type ParentLookupStack struct {
	arr []TargetLookup
}

//
func (this *ParentLookupStack) FindParent(gobj api.Instance) (ret api.Instance) {
	count := len(this.arr)
	if count > 0 {
		parentLookup := this.arr[count-1]
		ret = parentLookup(gobj)
	}
	return ret
}

//
func (this *ParentLookupStack) PushLookup(p TargetLookup) {
	this.arr = append(this.arr, p)
}

//
func (this *ParentLookupStack) PopLookup() {
	count := len(this.arr)
	this.arr = this.arr[0 : count-1]
}
