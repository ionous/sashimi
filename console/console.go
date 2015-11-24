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
	Close()
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
func (c SimpleConsole) Print(args ...interface{}) {
	fmt.Print(args...)
	fmt.Print(" ")
}

//
func (c SimpleConsole) Println(args ...interface{}) {
	fmt.Print(args...)
	fmt.Println()
}

//
func (c SimpleConsole) Readln() (ret string, okay bool) {
	okay = c.scanner != nil
	if okay {
		fmt.Print(">")
		okay = c.scanner.Scan()
		if okay {
			ret = c.scanner.Text()
		}
	}
	return ret, okay
}

//
func (c SimpleConsole) Close() {
}
