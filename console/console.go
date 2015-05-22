package console

import (
	"bufio"
	"fmt"
	"os"
)

//
type IConsole interface {
	Print(...interface{})
	Println(...interface{})
	Readln() (string, bool)
}

//
// for a main package test: create a console, and echo all input
func Echo() {
	c := NewConsole()
	for {
		if s, ok := c.Readln(); !ok {
			break
		} else {
			c.Println(s)
		}
	}
}

//
// Creates a SimpleConsole
func NewConsole() IConsole {
	scanner := bufio.NewScanner(os.Stdin)
	return SimpleConsole{scanner}
}

//
// implements IConsole
type SimpleConsole struct {
	scanner *bufio.Scanner
}

//
func (this SimpleConsole) Print(args ...interface{}) {
	fmt.Print(args...)
	fmt.Print(" ")
}

//
func (this SimpleConsole) Println(args ...interface{}) {
	fmt.Print(args...)
	fmt.Println()
}

//
func (this SimpleConsole) Readln() (ret string, okay bool) {
	okay = this.scanner != nil
	if okay {
		fmt.Print(">")
		okay = this.scanner.Scan()
		if okay {
			ret = this.scanner.Text()
		}
	}
	return ret, okay
}
