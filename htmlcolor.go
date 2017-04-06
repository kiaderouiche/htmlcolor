package htmlcolor

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

type SprintfFuncer interface {
	SprintfFunc() func(format string, a ...interface{}) string
}

type Formatter struct {
	TagColor     SprintfFuncer
	BracketColor SprintfFuncer
	CommentColor SprintfFuncer
	AttrKeyColor SprintfFuncer
	AttrValColor SprintfFuncer

	Indent bool
}

func NewFormatter(indent bool) *Formatter {
	return &Formatter{
		TagColor:     color.New(color.FgMagenta),
		BracketColor: color.New(color.FgCyan),
		CommentColor: color.New(color.FgBlue),
		AttrKeyColor: color.New(color.FgGreen),
		AttrValColor: color.New(color.FgRed),

		Indent: indent,
	}
}

func (f *Formatter) TagFprint(dst io.Writer, token html.Token) {
	sprintfAttrkey := f.AttrKeyColor.SprintfFunc()
	sprintfAttrval := f.AttrValColor.SprintfFunc()

	for _, attr := range token.Attr {
		fmt.Fprint(dst, " ", sprintfAttrkey(attr.Key+"="))
		fmt.Fprint(dst, sprintfAttrval("\""+attr.Val+"\""))
	}
}

func (f *Formatter) Format(dst io.Writer, src []byte) {
	sprintfTag := f.TagColor.SprintfFunc()
	sprintfBracket := f.BracketColor.SprintfFunc()
	sprintfComment := f.CommentColor.SprintfFunc()

	BracketOpen := "<"
	BracketClose := ">"

	srcReader := bytes.NewReader(src)
	tokenizer := html.NewTokenizer(srcReader)

	for {
		tokentype := tokenizer.Next()

		if tokentype == html.ErrorToken {
			return
		}

		token := tokenizer.Token()

		switch tokentype {
		case html.CommentToken:
			fmt.Fprint(dst, sprintfComment(token.String()))

		case html.DoctypeToken:
			fmt.Fprint(dst, sprintfComment(token.String()))

		case html.StartTagToken:
			fmt.Fprint(dst, sprintfBracket(BracketOpen))
			fmt.Fprint(dst, sprintfTag(token.Data))

			f.TagFprint(dst, token)

			fmt.Fprint(dst, sprintfBracket(BracketClose))

		case html.EndTagToken:
			fmt.Fprint(dst, sprintfBracket(BracketOpen+"/"))
			fmt.Fprint(dst, sprintfTag(token.Data))
			fmt.Fprint(dst, sprintfBracket(BracketClose))

		case html.SelfClosingTagToken:
			fmt.Fprint(dst, sprintfBracket(BracketOpen))
			fmt.Fprint(dst, sprintfTag(token.Data))

			f.TagFprint(dst, token)

			fmt.Fprint(dst, sprintfBracket(" /"+BracketClose))

		case html.TextToken:
			fmt.Fprint(dst, token.Data)
		}
	}
}
