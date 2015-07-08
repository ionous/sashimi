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
func (hs *History) Add(s string, mark HistoryMarker) HistoryMarker {
	hs.Restore(mark)
	if s != "" && hs.r.Value != s {
		hs.r.Value = s
		hs.r = hs.r.Next()
	}
	return nil
}

//
// Move backwards in time, returning the most recent item, then the next most recent, and so on.
// Returns `true` if there was any history to return.
//
func (hs *History) Back() (string, bool) {
	return hs._update(hs.r.Prev())
}

//
// Move forwards in time, returning the item one step closer to the most recent item.
// Returns `true` if not already at the most recent item.
//
func (hs *History) Forward() (string, bool) {
	return hs._update(hs.r.Next())
}

//
// Remember the current point in time. See: Restore()
//
func (hs *History) Mark() HistoryMarker {
	return hs.r
}

//
// Restore a previously remembered point in time.
//
func (hs *History) Restore(mark HistoryMarker) {
	if mark != nil {
		hs.r = mark
	}
}

//
// helper for forward and back
//
func (hs *History) _update(r *ring.Ring) (ret string, okay bool) {
	if s, ok := r.Value.(string); ok && s != "" {
		hs.r = r
		ret, okay = s, ok
	}
	return ret, okay
}
