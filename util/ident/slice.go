package ident

// IdSlice interface for sort
type IdSlice []Id

func (p IdSlice) Len() int           { return len(p) }
func (p IdSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IdSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
