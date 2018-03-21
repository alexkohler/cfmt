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
	write := flag.Bool("w", false, "write changes to file")
	flag.Usage = usage
	flag.Parse()

	if err := wrapComments(flag.Args(), *maxCommentLength, *write); err != nil {
		log.Println(err)
	}
}
