package main

import (
	"flag"
	"go/build"
	"log"
	"os"
)

const (
	pwd = "./"
)

func init() {
	build.Default.UseAllFiles = true
}

func usage() {
	log.Printf("Usage of %s:\n", os.Args[0])
	log.Printf("\ncfmt [flags] # runs on package in current directory\n")
	log.Printf("\ncfmt [flags] [packages]\n")
	log.Printf("Flags:\n")
	flag.PrintDefaults()
}

func main() {

	// Remove log timestamp
	log.SetFlags(0)

	maxCommentLength := flag.Uint("m", 80, "max comment length")
	dryRun := flag.Bool("dry", false, "dry run (print the changes that could be made, but don't modify any files)")
	verbose := flag.Bool("v", false, "print what was changed or ignored")
	flag.Usage = usage
	flag.Parse()

	if err := wrapComments(flag.Args(), *maxCommentLength, *dryRun, *verbose); err != nil {
		log.Println(err)
	}
}
