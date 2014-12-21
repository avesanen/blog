package main

import (
	"html/template"
)

type Article struct {
	Title    string        `json:"title"`
	Markdown string        `json:"markdown"`
	Content  template.HTML `json:"content"`
	Comments bool          `json:"comments"`
}
