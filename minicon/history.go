package minicon

import "container/ring"

//
// Keep a bunch of strings (eg. a user command buffer.)
//
type History struct {
	r *ring.Ring
}

//
// Records a stich in time.
//
type HistoryMarker *ring.Ring

//
// Create a new history buffer with the passed capacity.
//
func NewHistory(capacity int) *History {
	return &History{ring.New(capacity)}
}

//
// Add a new item to the current point in history.
// Takes an optional point in time which is first Restore()d
// Most people want to create a new most recent item, not overwrite something back in time.
//
func (this *History) Add(s string, mark HistoryMarker) HistoryMarker {
	this.Restore(mark)
	if s != "" && this.r.Value != s {
		this.r.Value = s
		this.r = this.r.Next()
	}
	return nil
}

//
// Move backwards in time, returning the most recent item, then the next most recent, and so on.
// Returns `true` if there was any history to return.
//
func (this *History) Back() (string, bool) {
	return this._update(this.r.Prev())
}

//
// Move forwards in time, returning the item one step closer to the most recent item.
// Returns `true` if not already at the most recent item.
//
func (this *History) Forward() (string, bool) {
	return this._update(this.r.Next())
}

//
// Remember the current point in time. See: Restore()
//
func (this *History) Mark() HistoryMarker {
	return this.r
}

//
// Restore a previously remembered point in time.
//
func (this *History) Restore(mark HistoryMarker) {
	if mark != nil {
		this.r = mark
	}
}

//
// helper for forward and back
//
func (this *History) _update(r *ring.Ring) (ret string, okay bool) {
	if s, ok := r.Value.(string); ok && s != "" {
		this.r = r
		ret, okay = s, ok
	}
	return ret, okay
}
