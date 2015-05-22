# Sashimi
An experimental IF engine written in Go.

Sashimi is a small platform to play with interactive fiction concepts. It is not yet a complete IF engine, and may never be. ( It's not an immediate goal, for instance, to provide proper English output. )

## Motivation
Based on some ideas from Infom7, Sashimi provides a "scripting language" embedded inside of Go. 

### Why Go?
[Go](http://golang.org) hits a nice middle-ground between scripting and c-like programming.

1. No semi-colons makes it read slightly more English like.
2. Back-quotes (``) for strings which makes embedding dialog much easier. (``"Much easier," she said.``)
3. Quick to compile.
4. The AST package opens the possibility of extracting scripts and transforming them into other languages for other runtimes. ( ex. lua, javascript, or possibly C# for Unity. )

### Why not Inform?
Inform is great, why not just experiment with it?

1. Inform is easy to read, but I find it's often difficult to write stories correctly; it always takes me a fair bit of trial and error to get what I want.
2. Not easily hackable, and I want to try new features.
3. More easily support stories split into multiple files for multiple authors.
4. More explicitly support 2D and 3D graphics.

## Status and Goals
Sashimi is currently capable of handling short, one-room stories; but, it lacks many standard IF features, including save-load.

Possible future features, in no particular order:
* A process for mocking up simple point-and-click adventure games -- in progress.
* A more complete set of IF features, especially: movement from room-to-room, doors, clothing, dialog.
* Documentation for the current architecture: especially event handling.
* Continuous time ( instead of purely turn-based ).
* Sqlite storage for the game world ( currently, the story must be re-compiled each time the game is run. )
* Save and load support.
* Improved scripting: type injection for script callbacks, global variables, type inferencing, variants for object properties, prettier syntax, ...
* Improved relationship support: especially relation-by-value; use sql(ite) to represent inter-object relationships in the runtime.
* Improved support for testing stories: especially a way to test expected output.
* Improved object modeling: for instance, [context state machines](https://github.com/ionous/hsm-statechart) to support concepts such as "lockable", "openable", "rideable" instead of classes.
* AST translation of script callbacks into other programming languages.
* Some sort of web-based story editor.

## Examples

A version of [A Day For Fresh Sushi](http://ifdb.tads.org/viewgame?id=7yiyxcnrlwejoffd) ported from Inform7 with permission.
'''
cd $GOPATH/src/github.com/ionous/sashimi/examples/fishy
go run fishy.go
(type "q" to quit.)
'''
Command line options:
* --text: use the simplier text console ( default is the fancier "minicon" which has a status bar, colors, etc. )
* --verbose: prints full output of actions and events as they happen.
* --dump: print all script generated classes,instances,actions, etc. to stdout, then exit.