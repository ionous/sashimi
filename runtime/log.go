package runtime

import "fmt"

//
type Log interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type LogAdapter struct {
	print func(s string)
}

func (log LogAdapter) Printf(format string, v ...interface{}) {
	log.print(fmt.Sprintf(format, v...))
}

func (log LogAdapter) Println(v ...interface{}) {
	log.print(fmt.Sprintln(v...))
}
