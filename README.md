# Sashimi
Sashimi provides a small platform to play with interactive fiction programming concepts. It is not yet a complete IF engine, and may never be. ( It's not an immediate goal, for instance, to provide proper English output. )

Based on some ideas from [Infom7](http://inform7.com), Sashimi uses a "scripting language" embedded inside of [Go](http://golang.org). 

For example:
```go
    s.The("container",
      Called("the cabinet"), In("the studio"),
      Is("openable", "closed").And("fixed in place"),
      Has("brief", "A huge cabinet, in the guise of an armoire, stands between the windows."),
      Has("description", "Large, and with a bit of an Art Nouveau theme going on in the shape of the doors."),
    )
    s.The("cabinet", IsKnownAs("armoire"))
    lookedUnderCabinet := false 
    s.The("cabinet", 
        When("looking under").Always(func(g G.Play) {
          if !lookedUnderCabinet {
              lookedUnderCabinet = true
              g.The("evil fish").Says(`"Dustbunnies," predicts the fish, with telling accuracy.`)
              g.StopHere()
          }
      }))
```

## Motivation

### Why Go?
Go hits a nice middle-ground between scripting and C-like programming.

1. No semi-colons makes it slightly more English like than other options.
2. Back-quotes ``(`)`` for strings makes embedding dialog much easier. (`` `"Much easier," she said.` ``)
3. Quick to compile.
4. The AST package opens the possibility of extracting script callbacks and transforming them into other languages for other runtimes ( ex. lua with a custom C-runtime, javascript in a web-app, or C# for Unity. )

### Why not Inform?
Inform is great, why not just experiment with it?

1. Inform is easy to read, but I find it's often difficult to write stories correctly; it always takes me a fair bit of trial and error to get what I want.
2. Not easily hackable, and I want to try new low-level features.
3. I would like to more easily split a single story into multiple files (for instance: for multiple authors.)
4. I would like more explicit support for 2D and 3D graphics.

## Status and Goals
Sashimi is currently capable of handling short, one-room stories; but, it lacks many standard IF features, including save-load.

Possible future features, in no particular order:
* A process for mocking up simple point-and-click adventure games ( *in progress*. )
* A more complete set of IF features, especially: movement from room-to-room, doors, clothing, dialog.
* Documentation for the current architecture: especially event handling.
* Continuous time ( instead of purely turn-based. )
* Sqlite storage for the game world, with the ability to merge data about story generated objects from non-story sources.
* Save and load support.
* Improved scripting: type injection for script callbacks, global variables, type inferencing, variants for object properties, prettier syntax, ...
* Improved relationship support: especially relation-by-value; use sql(ite) to represent inter-object relationships in the runtime.
* Improved support for testing stories: especially a way to test expected output.
* Improved object modeling: for instance, [context state machines](https://github.com/ionous/hsm-statechart) to support concepts such as "lockable", "openable", "rideable" ( instead of classes. )
* AST translation of script callbacks into other programming languages.
* Some sort of web-based story editor.

## Examples

A version of [A Day For Fresh Sushi](http://ifdb.tads.org/viewgame?id=7yiyxcnrlwejoffd) ported from Inform7 with permission.

```
cd $GOPATH/src/github.com/ionous/sashimi/examples/fishy
go run fishy.go
(type "q" to quit.)
```

Command line options:
&nbsp;&nbsp;&nbsp;--text: use the simplier text console ( default is the fancier "minicon" which has a status bar, colors, etc. )
&nbsp;&nbsp;&nbsp;--verbose: prints full output of actions and events as they happen.
&nbsp;&nbsp;&nbsp;--dump: print all script generated classes,instances,actions, etc. to stdout, then exit.