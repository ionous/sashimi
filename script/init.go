package script

type InitCallback func(s *Script)

var allInits []InitCallback

// AddScript supports callbacks of this pattern:
/*
 func init() {
	AddScript(func(s *Script) {
		s.The...
	})
	}
*/
func AddScript(cb InitCallback) {
	allInits = append(allInits, cb)
}

func InitScripts() *Script {
	s := &Script{}
	for _, initCb := range allInits {
		initCb(s)
	}
	return s
}
