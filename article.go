package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
)

type Article struct {
	Id       string        `json:"id"`
	Title    string        `json:"title"`
	Markdown string        `json:"markdown"`
	Content  template.HTML `json:"content"`
	Comments bool          `json:"comments"`
}

func (a *Article) render() {
	unsafe := blackfriday.MarkdownCommon([]byte(a.Markdown))
	a.Content = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}
