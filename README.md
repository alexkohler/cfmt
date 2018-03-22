# cfmt

cfmt is a tool to wrap Go comments over a certain length to a new line.

## Installation

    go get -u github.com/alexkohler/cfmt

    **Note**: cfmt requires gofmt. 

## Usage

Similar to other Go static anaylsis tools (such as golint, go vet), prealloc can be invoked with one or more filenames, directories, or packages named by its import path. Prealloc also supports the `...` wildcard. 

    cfmt [flags] files/directories/packages

### Flags
- **-w** - Writes changes to file. By default, `cfmt` will only print the changes it will make, but will not modify the input files.

## Example

`cfmt` will wrap to a new line or join an existing line as appropriate. For example, running `cfmt -m=100` on the following file:

```Go
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

```

After running `cfmt -w`:

```Go
// I am a long comment that is over 80 characters long. I should probably wrap
// to a new line.
func test() {
	// I am a long comment that is over 80 characters long. I should probably wrap
	// below to the rest of the comment.
	fmt.Println("hello world")

	// I am a long comment that is over 80 characters long. I should probably wrap
	// below to the rest of the comment.

	//I am a long comment that starts without a space and is over 80 characters
	//long. When I wrap, I should still start without a space

	/* I am a block comment. I get ignored. I can be waaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaay longer than 80 characters and still won't be affected.
	 */

	// Below is a long, single word comment. This should be ignored because it's
	// usually indicative of a diagram, divider, etc...
	// -*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*
}

```

`cfmt` ignores block comments.


## Contributing

Pull requests welcome!


## Other static analysis tools

If you've enjoyed prealloc, take a look at my other static anaylsis tools!
- [nakedret](https://github.com/alexkohler/nakedret) - Finds naked returns.
- [unimport](https://github.com/alexkohler/unimport) - Finds unnecessary import aliases.
- [prealloc](https://github.com/alexkohler/prealloc) - Finds slice declarations that could potentially be preallocated.