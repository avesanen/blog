package main

import (
	"github.com/zenazn/goji/web"
	"html/template"
	"log"
	"net/http"
	"time"
)

// /:archive/:chapter/:page

type Archive struct {
	Id       string
	Title    string   `json:"title"`
	Chapters []string `json:"chapters"`
}

func (a *Archive) getChapter(id string) (*Chapter, error) {
	var c Chapter
	if err := db.Read("chapters", id, &c); err == nil {
		return &c, nil
	} else {
		return nil, err
	}
}

func (a *Archive) getChapters() ([]*Chapter, error) {
	chapters := make([]*Chapter, 0)
	for _, v := range a.Chapters {
		if c, err := a.getChapter(v); err == nil {
			chapters = append(chapters, c)
		} else {
			log.Println(err.Error())
			return nil, err
		}
	}
	return chapters, nil
}

type Chapter struct {
	Id    string
	Title string   `json:"title"`
	Pages []string `json:"pages"`
}

type Page struct {
	Id          string
	Title       string        `json:"title"`
	Markdown    string        `json:"markdown"`
	Content     template.HTML `json:"content"`
	Comments    bool          `json:"comments"`
	PublishDate time.Time     `json:"time"`
}

func readPage() *Page {
	return nil
}

func ListArchivesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func ViewArchiveHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	archiveId := c.URLParams["archiveId"]
	a := map[string]interface{}{}
	var v Archive
	if err := db.Read("archive", archiveId, &v); err == nil {
		a["Archive"] = v
		if chapters, err := v.getChapters(); err == nil {
			a["Chapters"] = chapters
		}
		log.Println(a)
		renderTemplate(w, r, "view_archive", a)
	} else {
		http.Error(w, "not found", http.StatusNotFound)
	}
}
