package api

//
type Log interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}
