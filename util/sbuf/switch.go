package sbuf

//return sbuf.New(n.Id, sbuf.Q{'(', n.Class, ')'}).String()
// type Q struct {
// 	Middle SBuf // ""
// 	}
// rather than String() for each prom, we will want Write(buf, sep)
// eval parms on New, but only write on outer String()
func New(parts ...interface{}) Switch {
	return Switch(parts)
}

type Switch []interface{}

func (j Switch) Append(i interface{}) Switch {
	return append(j, i)
}
func (j Switch) String() string {
	return j.Join("")
}
func (j Switch) Join(sep string) string {
	// ugly slow. fix.
	parts := make([]Stringer, len(j))
	for i, p := range j {
		var str Stringer
		switch v := p.(type) {
		case bool:
			str = Bool{v}
		case int:
			str = Int{v}
			// rune and int32 are the same type!?!
		// case rune:
		// 	str = Rune{v}
		case int32:
			str = Int64{int64(v)}
		case int64:
			str = Int64{v}
		case uint:
			str = Hex64{uint64(v)}
		case uint32:
			str = Hex64{uint64(v)}
		case uint64:
			str = Hex64{v}
		case float32:
			str = Float{float64(v)}
		case float64:
			str = Float{v}
		case Stringer:
			str = v
		case error:
			str = Error{v}
		case string:
			str = String{v}
		default:
			str = String{"###"}
		}
		parts[i] = str
	}
	return (&StringBuffer{parts}).Join(" ")
}
