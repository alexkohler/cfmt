package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func Test_processComments(t *testing.T) {
	type args struct {
		src              string
		maxCommentLength int
	}
	tests := []struct {
		name         string
		args         args
		wantComments [][]string
	}{
		{name: "single line long block comment",
			args: args{
				src: `
				  package main

				  /* i am longer than 5 characters */
				  func main() {
				  	println("Hello, World!")
				  }
				  `,
				maxCommentLength: 5,
			},
			wantComments: [][]string{
				{
					// This comment be ignored (it's a block comment)
					"/* i am longer than 5 characters */",
				},
			},
		},
		{name: "second line long block comment",
			args: args{
				src: `
				  package main

				  /* short
				  i am longer than 10 characters */
				  func main() {
				  	println("Hello, World!")
				  }
				  `,
				maxCommentLength: 10,
			},
			wantComments: [][]string{
				{
					// This comment be ignored (it's a block comment)
					`/* short
				  i am longer than 10 characters */`,
				},
			},
		},
		{name: "monolithic comment with separate wrappable group",
			args: args{
				src: `
				  package main
				  
				  //1234567891011121314
				  func main() {
					// i should still be fixed
				  	println("Hello, World!")
				  }
				  `,
				maxCommentLength: 12,
			},
			wantComments: [][]string{
				{
					// text in a long conglomerate string generally hints towards a
					// diagram that shouldn't be formatted
					"//1234567891011121314",
				},
				{
					"// i should",
					"// still be fixed",
				},
			},
		},
		{name: "comment with another comment to wrap to",
			args: args{
				src: `
				  package main

				  // i am long long long long long
				  // i am short
				  func main() {
				  	println("Hello, World!")
				  }
				  `,
				maxCommentLength: 25,
			},
			wantComments: [][]string{
				{
					// Comment should wrap
					"// i am long long long",
					"// long long i am short",
				},
			},
		},
		{name: "single line comment with no beginning space",
			args: args{
				src: `
				  package main

				  //i am a long comment with no beginning space
				  func main() {
				  	println("Hello, World!")
				  }
				  `,
				maxCommentLength: 25,
			},
			wantComments: [][]string{
				{
					// This comment be ignored (it's a block comment)
					"//i am a long comment",
					"//with no beginning space",
				},
			},
		},
		{name: "large comment group (ignored)",
			args: args{
				src: `
				  package main

				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  //i am a long comment with no beginning space
				  func main() {
				  	println("Hello, World!")
				  }
				  `,
				maxCommentLength: 25,
			},
			wantComments: [][]string{
				{
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
					"//i am a long comment with no beginning space",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fset := token.NewFileSet() // positions are relative to fset
			f, err := parser.ParseFile(fset, "", tt.args.src, parser.ParseComments)
			if err != nil {
				t.Fatalf("Did not expect error parsing file, %v", err)
			}

			processComments(fset, f.Comments, tt.args.maxCommentLength, false)

			if len(tt.wantComments) != len(f.Comments) {
				t.Fatalf("Expected group length %v, got %v", len(tt.wantComments), len(f.Comments))
			}

			for j, commentGroup := range f.Comments {

				if len(tt.wantComments[j]) != len(commentGroup.List) {
					t.Fatalf("Expected comment group length %v, got %v", len(tt.wantComments[j]), len(commentGroup.List))
				}

				for i, actual := range commentGroup.List {
					if tt.wantComments[j][i] != actual.Text {
						t.Errorf("Mismatched comments.\nwant %v\n got %v\n", tt.wantComments[j][i], actual.Text)
					}
				}
			}

		})
	}
}
