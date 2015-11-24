package extract

import (
	"github.com/ionous/sashimi/util/errutil"
	"go/ast"
	"go/parser"
	"go/token"
)

type Snip func(file string, line int, snippet []byte) error

// ParseSource, the original .go files on disk, and generate file, line, and source code snippers for every callback encountered.
func ParseSource(file string, bytes []byte, snip Snip) (err error) {
	fset := token.NewFileSet()
	// filename is a label for src is not nil ( and is: string, []byte, or io.Reader)
	if f, e := parser.ParseFile(fset, file, bytes, 0); e != nil {
		err = e
	} else {
		literally := visitorStruct{}
		literally.visit = func(n ast.Node) (w ast.Visitor) {
			if literal, ok := n.(*ast.FuncLit); ok {
				rets, params := literal.Type.Results, literal.Type.Params
				// there probably are other G.Play callbacks... if they were exported, we could link to them automatically.
				if rets == nil && len(params.List) == 1 {
					g := params.List[0]
					// found by print %T experimentation
					if s, ok := g.Type.(*ast.SelectorExpr); ok {
						if s.Sel.Name == "Play" {
							start, end := fset.Position(n.Pos()), fset.Position(n.End())
							sub := bytes[start.Offset:end.Offset]
							//blockStart := fset.Position(block.Pos()).Line
							if e := snip(file, start.Line, sub); e != nil {
								err = errutil.Append(err, e)
							}
							return nil // dont go deeper than the literal
						}
					}
				}
			}
			// keep opening up to find (other) literals
			if err != nil {
				return nil
			} else {
				return literally
			}
		} // literally.visit
		ast.Walk(literally, f)
	}
	return err
}

type visitorStruct struct {
	visit func(ast.Node) (w ast.Visitor)
}

func (v visitorStruct) Visit(n ast.Node) (w ast.Visitor) {
	return v.visit(n)
}
