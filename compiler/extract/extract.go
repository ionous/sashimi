package extract

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type CB func(file string, line int, bytes []byte)

func Extract(file string, bytes []byte, cb CB) (err error) {
	fset := token.NewFileSet()
	// filename is a label for src is not nil ( and is: string, []byte, or io.Reader)
	if f, e := parser.ParseFile(fset, file, bytes, 0); e != nil {
		err = e
	} else {
		find := VisitorStruct{}
		find.visit = func(n ast.Node) (w ast.Visitor) {
			if _, ok := n.(*ast.FuncDecl); ok {
				// scan for the function block:
				blocker := VisitorStruct{}
				blocker.visit = func(n ast.Node) (w ast.Visitor) {
					// we get three: Indent, FuncType, and BlockStmt
					if block, ok := n.(*ast.BlockStmt); ok {
						// scan for function literals
						literally := VisitorStruct{}
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
											blockStart := fset.Position(block.Pos()).Line
											cb(file, blockStart, sub)
										}
									}
								}
								return nil
							}
							return literally
						} // literally.visit
						return literally
					}
					return nil
				}
				return blocker
			}
			return find
		}
		ast.Walk(find, f)
	}
	return err
}

type VisitorStruct struct {
	visit func(ast.Node) (w ast.Visitor)
}

func (v VisitorStruct) Visit(n ast.Node) (w ast.Visitor) {
	return v.visit(n)
}
