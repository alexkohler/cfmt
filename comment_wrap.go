package main

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

func wrapComments(args []string, maxCommentLength uint, dryRun, verbose bool) error {

	fset := token.NewFileSet()

	files, err := parseInput(args, fset)
	if err != nil {
		return fmt.Errorf("could not parse input %v", err)
	}

	for _, f := range files {
		processComments(fset, f.Comments, int(maxCommentLength), dryRun, verbose)

		file, err := os.Create(fset.File(f.Pos()).Name())
		if err != nil {
			return err
		}
		defer file.Close()
		if err := printer.Fprint(file, fset, f); err != nil {
			log.Fatal(err)
		}

	}

	return nil
}

func processComments(fset *token.FileSet, comments []*ast.CommentGroup, maxCommentLength int, dryRun, verbose bool) {

	for _, cg := range comments {
		cIdx := 0
		for _, c := range cg.List {

			if len(c.Text) > maxCommentLength {

				file := fset.File(cg.Pos())
				lineNumber := file.Position(cg.List[cIdx].Pos()).Line

				// Block comments are ignored for now to simplify logic
				if strings.Contains(c.Text, "/*") || strings.Contains(c.Text, "*/") {
					if verbose {
						log.Printf("%v:%v ignoring block comment\n", file.Name(), lineNumber)
					}
					break
				}

				// Split on each word
				words := strings.Split(c.Text, " ")
				// This means we have one giant word longer than X characters. It's probably a diagram or something, so we'll leave it.
				if len(words) == 1 {
					continue
				}

				// Otherwise, chop off words in reverse and glob them into a new line
				currentLength := len(c.Text)
				var choppedWords []string
				for i := len(words) - 1; i >= 0; i-- {
					currentLength = currentLength - len(words[i])
					currentLength-- // subtract another 1 because we split on a space
					choppedWords = append(choppedWords, words[i])
					// Chopping off this word fixed our problems

					if currentLength < maxCommentLength {
						cg.List[cIdx].Text = strings.Join(words[0:i], " ")
						// See if this is in a comment group we can wrap below to
						if cIdx < len(cg.List)-1 {
							// we were going in reverse, so append these in the other order.
							reverse(choppedWords)

							// Split on the comment
							splitComment := strings.Split(cg.List[cIdx+1].Text, "//")

							// See whether or not this comment follows the "// comment" or "// comment" idiom
							commentStr := "//"
							usesSpace := strings.HasPrefix(cg.List[cIdx+1].Text, "// ")
							if usesSpace {
								commentStr += " "
							}

							cg.List[cIdx+1].Text = commentStr + strings.Join(choppedWords, " ") + splitComment[0]

							// Tack the rest of the comment back on
							if len(splitComment) > 1 {
								for _, s := range splitComment {
									cg.List[cIdx+1].Text += s
								}
							}
							break
						} else {
							// Otherwise, we have no room to wrap. This is easy, just create a new comment
							reverse(choppedWords)

							// See whether or not the preceding comment follows the "// comment" or "//comment" idiom
							commentStr := "//"
							usesSpace := strings.HasPrefix(cg.List[cIdx].Text, "// ")
							if usesSpace {
								commentStr += " "
							}

							cg.List = append(cg.List, &ast.Comment{Slash: 0, Text: commentStr + strings.Join(choppedWords, " ")})
							break
						}
					}

				}

				if dryRun || verbose {
					log.Printf("%v:%v can be reduced to\n", file.Name(), lineNumber)
					log.Printf("    %v\n", cg.List[cIdx].Text)
					log.Printf("    %v\n", cg.List[cIdx+1].Text)
				}

			}

			cIdx++
		}
	}
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
