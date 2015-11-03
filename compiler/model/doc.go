/*
Model provides a post-compiled representation of the scripts.
( It should probably be assumed to be internally consistent. )
The table version of things would also build,merge into this same code.
The script callbacks currently must always be re-compiled.
The idea is that they could be extracted somehow via the ast, and perhaps recompiled as
( And if not, probably not the end of the world -- new data could still be merged into the model from tables,
the original scripts just couldnt be serialized out to non-go files. )
*/
package model
