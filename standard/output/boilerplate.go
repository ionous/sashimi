package output

import (
	"fmt"
	"github.com/ionous/mars/script/backend"
	"github.com/ionous/mars/std"
	"github.com/ionous/sashimi/compiler"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/play"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/standard/framework"
	"io"
	"io/ioutil"
)

// simplest interface:
func Run(s backend.Declaration) {
	RunGame(s, ParseCommandLine())
}

func GetWriter(opt Options, w io.Writer) (ret io.Writer) {
	if opt.verbose {
		ret = w
	} else {
		ret = ioutil.Discard
	}
	return
}

func GetConsole(opt Options) (ret C.IConsole) {
	if opt.hasConsole {
		ret = opt.cons
	} else {
		ret = C.NewConsole()
	}
	return
}

func RunGame(s backend.Declaration, opt Options) (err error) {
	cons := GetConsole(opt)
	debugWriter := GetWriter(opt, cons)
	src := S.Statements{}
	if e := std.Std().Generate(&src); e != nil {
		err = e
	} else if e := s.Generate(&src); e != nil {
		err = e
	} else if model, e := compiler.Compile(src, debugWriter); e != nil {
		err = e
	} else if opt.dump {
		// var b bytes.Buffer
		// enc := gob.NewEncoder(&b)
		// gob.Register(ident.Id(""))
		// if e := enc.Encode(model.Model); e != nil {
		// 	panic(e)
		// }
		// fmt.Println(fmt.Sprintf("size: %d(b) %.2f(k)", b.Len(), float64(b.Len())/1024.0))
		// model.Model.PrintModel(func(args ...interface{}) { fmt.Println(args...) })
		panic("not implemented")
	} else {
		vals := make(metal.ObjectValueMap)
		modelApi := metal.NewMetal(model.Model, vals)
		parents := framework.NewParentLookup()
		g := play.
			NewConfig().
			SetWriter(cons).
			SetLogger(debugWriter).
			SetParentLookup(parents).
			MakeGame(modelApi)
		parents.Run = g
		err = PlayGame(cons, g)
	}
	if err != nil {
		panic(err)
	}
	return
}

func PlayGame(cons C.IConsole, g play.Game) (err error) {
	if game, e := framework.NewStandardGame(g); e != nil {
		err = e
	} else if game, e := game.Start(); e != nil {
		err = e
	} else {
		for {
			if s, ok := cons.Readln(); !ok {
				break
			} else if e := game.Input(s); e != nil {
				fmt.Fprintln(cons, e)
			} else if game.IsQuit() || game.IsComplete() {
				break
			}
		}
	}

	return err
}
