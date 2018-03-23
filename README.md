# cfmt

cfmt is a tool to wrap Go comments over a certain length to a new line.

## Installation

`go get -u github.com/alexkohler/cfmt`

**Note**: cfmt requires gofmt. 

## Usage

Similar to other Go static anaylsis tools (such as golint, go vet), cfmt can be invoked with one or more filenames, directories, or packages named by its import path. cfmt also supports the `...` wildcard. 

    cfmt [flags] files/directories/packages

### Flags
- **-m** - Maximum comment length before cfmt should insert a line break.  (Default 80)
- **-w** - Writes changes to file. By default, `cfmt` will only print the changes it will make, but will not modify the input files.

## Examples

`cfmt` will wrap to a new line or join an existing line as appropriate. See the following before/afters of running `cfmt -m=100`:

**Before**
```Go
// I am a long comment that is over 100 characters long. I should probably wrap to a new line.
```

**After**
```Go
// I am a long comment that is over 100 characters long. I should probably wrap
// to a new line.
```
---

**Before**
```Go
// I am a long comment that is over 100 characters long. I should probably wrap below to the
// rest of the comment.
```

**After**
```Go
// I am a long comment that is over 100 characters long. I should probably wrap
// below to the rest of the comment.
```
---

**Before**
```Go
//I am a long comment that starts without a space and is over 100 characters long. When I wrap, I should still start without a space
```

**After**
```Go
//I am a long comment that starts without a space and is over 100 characters
//long. When I wrap, I should still start without a space
```
---

`cfmt` ignores block (`/* */`) comments and "grouped" comments over a length of 10 (i.e. 10+ consecutive lines starting with `//`).


## Contributing

Pull requests welcome!


## Other static analysis tools

If you've enjoyed cfmt, take a look at my other static anaylsis tools!
- [nakedret](https://github.com/alexkohler/nakedret) - Finds naked returns.
- [unimport](https://github.com/alexkohler/unimport) - Finds unnecessary import aliases.
- [prealloc](https://github.com/alexkohler/prealloc) - Finds slice declarations that could potentially be preallocated.
