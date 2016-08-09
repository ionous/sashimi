package main

import "github.com/ionous/sashimi/minicon"

func main() {
	d := minicon.NewMiniCon()
	defer d.Close()
	d.Status.Left = "Console Test!"
	d.Status.Right = "Right text should run under the left text when the line is too long."

	d.Println("Type `q` to quit; press `esc` to speed up text.")
	d.Println("")
	// "github.com/drhodes/golorem" // lorem
	// for i := 0; i < 100; i++ {
	// 	str := lorem.Sentence(5, 100)
	// 	d.Println(str)
	// }
	d.Println("You're on the run. You've got a million errands to do -- your apartment to get cleaned up, the fish to feed, lingerie to buy, Britney's shuttle to meet--")
	d.Println("across town from anywhere else you have to do. Oh well, you'll")
	d.Println("Type `q` to quit; press `esc` to speed up text.")

	for {
		str := d.Update()
		if str == "q" {
			break
		}
		if str != "" {
			d.Println(str)
		}
	}
}
