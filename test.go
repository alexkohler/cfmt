package main

import "fmt"

// I am a long comment that is over 80 characters long. I should probably wrap to a new line.
func test() {
	// I am a long comment that is over 80 characters long. I should probably wrap below to the
	// rest of the comment.
	fmt.Println("hello world")

	// I am a long comment that is over 80 characters long. I should probably wrap below to the
	// rest of the comment.

	//I am a long comment that starts without a space and is over 80 characters long. When I wrap, I should still start without a space

	/* I am a block comment. I get ignored. I can be waaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaay longer than 80 characters and still won't be affected.
	 */

	// Below is a long, single word comment. This should be ignored because it's usually indicative of a diagram, divider, etc...
	// -*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*
}
