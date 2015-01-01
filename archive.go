package main

import (
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

/*
// /:archive/:chapter/:page

type Archive struct {
	Id       string
	Title    string   `json:"title"`
	Chapters []string `json:"chapters"`
}


func (a *Archive) addChapter()
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

func (a *Archive) getPage(pageNumber int, pageSize int) []*Chapter {
	start := pageNumber * pageSize
	end := start + pageSize

	if start > len(a.Chapters) {
		return nil
	}
	if end > len(a.Chapters) {
		end = len(a.Chapters) - start
	}
	return a.Chapters[start:end]
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
*/
