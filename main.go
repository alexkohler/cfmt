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
	log.Printf("\ncmt [flags] # runs on package in current directory\n")
	log.Printf("\ncmt [flags] [packages]\n")
	log.Printf("Flags:\n")
	flag.PrintDefaults()
}

func main() {

	// Remove log timestamp
	log.SetFlags(0)

	maxCommentLength := flag.Uint("m", 80, "max comment length")
	dryRun := flag.Bool("dry-run", false, "print what could potentially be changed")
	verbose := flag.Bool("v", false, "print what was changed")
	flag.Usage = usage
	flag.Parse()

	if err := wrapComments(flag.Args(), *maxCommentLength, *dryRun, *verbose); err != nil {
		log.Println(err)
	}
}
