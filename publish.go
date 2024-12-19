package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/shurcooL/github_flavored_markdown"
)

// top-level func which can converts markdown file to a pdf file
func main() {
	// Create a memory DOM from markdown using the blackfriday library
	// https://github.com/russross/blackfriday
	markdown := []byte(`
# Hello World

This is a **bold** text and this is *italic*.

## List Example
- Item 1
- Item 2
- Item 3

[Link to Google](https://www.google.com)
    `)

	// Convert markdown to HTML
	//    html := blackfriday.Run(markdown)

	// Print the resulting HTML
	//    fmt.Println(string(html))

	//html := blackfriday.R(markdown)

	// Alternative: Use goldmark library
	// https://github.com/yuin/goldmark
	// var buf bytes.Buffer
	// if err := goldmark.Convert(markdown, &buf); err != nil {
	// 	panic(err)
	// }
	// htmlOutput := buf.Bytes()

	// Alternative: Use github-markdown package
	// https://github.com/shurcooL/github_flavored_markdown
	html := github_flavored_markdown.Markdown(markdown)

	// Convert the HTML to a PDF file
	// Use chromedp package to convert HTML to PDF
	// First, create a new Chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+string(html)),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	); err != nil {
		log.Fatal(err)
	}

	// Write PDF to file
	if err := ioutil.WriteFile("output.pdf", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
