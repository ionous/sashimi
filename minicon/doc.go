/*
 Simple interactive fiction console based on http://github.com/nsf/termbox-go.

 Example usage:

	con := NewMiniCon() // create console
	defer con.Close()  // defer cleanup

	// print some example text
	con.Status.Left, con.Status.Right = "Left!", "Right!"
	con.Println("Type `q` to quit; press `esc` to speed up text.")

	// loop till done
	for {
		// read user input
		userInput := con.Update()
		if userInput == "q" {
			break
		}
		// echo to the display
		d.Println(userInput)
	}
*/
package minicon
