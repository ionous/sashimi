// Package appengine is an attempt at running sashimi on google's appengine
// It works, but is slow. Long term, removing the startup parsing: instead an pipeline to compile into data; and moving to proper queries would help.
// example:
// run gen.go in sash-alice/gen to build/extract the source
// run "goapp serve demo" in sash-alice
// goapp serve demo.
// http://localhost:8000/datastore
package appengine
