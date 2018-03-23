package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func wrapComments(args []string, maxCommentLength uint, write bool) error {

	fset := token.NewFileSet()

	files, err := parseInput(args, fset)
	if err != nil {
		return fmt.Errorf("could not parse input %v", err)
	}

	var wg sync.WaitGroup
	for _, f := range files {
		processComments(fset, f.Comments, int(maxCommentLength), write)

		wg.Add(1)

		go func(f *ast.File) {
			defer wg.Done()
			if write {
				// Write the changes out to the file
				fileName := fset.File(f.Pos()).Name()

				file, err := os.Create(fileName)
				if err != nil {
					log.Println(err)
					return
				}
				defer file.Close()

				if len(f.Comments) == 0 {
					return
				}

				if err := printer.Fprint(file, fset, f); err != nil {
					log.Fatal(err)
					return
				}

				cmd := exec.Command("gofmt", "-w", fileName)
				var out bytes.Buffer
				cmd.Stdout = &out
				err = cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		}(f)

	}

	wg.Wait()
	return nil
}

func processComments(fset *token.FileSet, comments []*ast.CommentGroup, maxCommentLength int, write bool) {

	for _, cg := range comments {
		cIdx := 0

		// Ignore huge block comment groups
		if len(cg.List) > 10 {
			continue
		}
		for _, c := range cg.List {

			if len(c.Text) > maxCommentLength {

				file := fset.File(cg.Pos())
				lineNumber := file.Position(cg.List[cIdx].Pos()).Line

				// Block comments are ignored
				if strings.HasPrefix(c.Text, "/*") || strings.HasSuffix(c.Text, "*/") {
					if !write {
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

							cg.List = append(cg.List, &ast.Comment{Slash: cg.End() + 1, Text: commentStr + strings.Join(choppedWords, " ") + "\n"})
							break
						}
					}

				}

				if !write {
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
