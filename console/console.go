package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

//
type IConsole interface {
	io.Writer
	Readln() (string, bool)
}

// for a main package test: create a console, and echo all input
func Echo() {
	c := NewConsole()
	for {
		if s, ok := c.Readln(); !ok {
			break
		} else {
			fmt.Fprintln(c, s)
		}
	}
}

// Creates a SimpleConsole
func NewConsole() IConsole {
	scanner := bufio.NewScanner(os.Stdin)
	return SimpleConsole{os.Stdout, scanner}
}

// implements IConsole
type SimpleConsole struct {
	io.Writer
	scanner *bufio.Scanner
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
	return
}
