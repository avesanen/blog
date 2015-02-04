package main

import (
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"sort"
	"time"
)

func NewArchive() *Archive {
	a := &Archive{}
	a.Id = randSeq(16)
	a.Pages = make([]*Page, 0)
	return a
}

type Archive struct {
	Id    string `json:"id"`
	Pages Pages  `json:"pages"`
}

func (a *Archive) AddPage(p *Page) {
	a.Pages = append(a.Pages, p)
}

func (a *Archive) RmPage(p *Page) {
	for i, v := range a.Pages {
		if v == p {
			a.Pages = append(a.Pages[:i], a.Pages[i+1:]...)
			return
		}
	}
}

func (a *Archive) SortByDate() {
	sort.Sort(byDate{a.Pages})
}

func NewPage() *Page {
	p := &Page{}
	p.Id = randSeq(16)
	p.PublishDate = time.Now()
	return p
}

type Page struct {
	Id          string    `json:"id"`
	PublishDate time.Time `json:"time"`
}

type Pages []*Page

func (a Pages) Len() int      { return len(a) }
func (a Pages) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type byDate struct{ Pages }

func (a byDate) Less(i, j int) bool { return a.Pages[i].PublishDate.After(a.Pages[j].PublishDate) }

func ViewArchiveHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	archiveId := c.URLParams["archiveId"]
	a := map[string]interface{}{}
	a["archiveId"] = archiveId
	page := r.URL.Query().Get("page")
	show := r.URL.Query().Get("show")
	log.Printf("Archive %v, page %v, show %v", archiveId, page, show)
	renderTemplate(w, r, "view_archive", a)
}

/*
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
*/
