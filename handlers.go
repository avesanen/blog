package main

import (
	"encoding/json"
	"fmt"
	"github.com/zenazn/goji/web"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func requiresLogin(h web.HandlerFunc) web.HandlerFunc {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		if s := getSession(r); s == nil {
			renderTemplate(w, r, "login_page", map[string]interface{}{})
		} else {
			h.ServeHTTPC(c, w, r)
		}
	}
	return fn
}

func renderTemplate(w http.ResponseWriter, r *http.Request, t string, a map[string]interface{}) {
	templates := template.Must(template.ParseGlob("./templates/*.tmpl"))
	a["Session"] = getSession(r)
	if err := templates.ExecuteTemplate(w, t, a); err != nil {
		log.Println("Error rendering template:", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/a/index", http.StatusSeeOther)
}

func articleViewHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	articleId := c.URLParams["article"]
	var article Article
	err := db.Read("article", articleId, &article)
	if err != nil {
		http.Redirect(w, r, "/a/"+articleId+"/edit", http.StatusFound)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"status": "ok", "url": "/" + articleId})
		return
	}
	renderTemplate(w, r, "view_page", map[string]interface{}{"Article": article})
}

func articleEditHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	articleId := c.URLParams["article"]
	if r.Method == "GET" {
		var article Article
		err := db.Read("article", articleId, &article)
		if err != nil {
			article = Article{}
			article.Id = articleId
			article.Title = articleId
			article.Content = ""
			article.Comments = false
			article.Markdown = ""
		}
		renderTemplate(w, r, "edit_page", map[string]interface{}{"Article": article})
	} else if r.Method == "POST" {
		var article Article
		err := db.Read("article", articleId, &article)
		if err != nil {
			article = Article{}
			article.Id = articleId
			article.Title = articleId
			article.Content = ""
			article.Comments = false
			article.Markdown = ""
			if err := db.Write("article", articleId, &article); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body:", err.Error())
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": err.Error()})
			http.Error(w, err.Error(), 500)
			return
		}

		form := make(map[string]string)

		if err = json.Unmarshal(body, &form); err != nil {
			log.Println("Error unmarshaling body:", err.Error())
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": err.Error()})
			http.Error(w, err.Error(), 500)
			return
		}

		article.Markdown = string(form["markdown"])
		article.Title = string(form["title"])
		if form["comments"] == "true" {
			article.Comments = true
		} else {
			article.Comments = false
		}
		article.render()

		if err := db.Write("article", articleId, &article); err != nil {
			log.Println("Error saving article:", err.Error())
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": err.Error()})
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"status": "ok", "url": "/a/" + articleId})
	}
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func uploadPostHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	formdata := r.MultipartForm
	files := formdata.File["files"]

	urls := make([]string, 0)

	for i, _ := range files {
		// Open sent file
		file, err := files[i].Open()
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Generate file struct
		ft := strings.SplitN(files[i].Header.Get("Content-Type"), "/", 2)
		if len(ft) != 2 {
			log.Println("Content-Type not in Type/Subtype format.")
			return
		}
		fType := ft[0]
		fSubType := ft[1]

		filename := files[i].Filename

		fSplit := strings.Split(filename, ".")
		fileExtension := fSplit[len(fSplit)-1]
		fileId := randSeq(16)

		filename = fileId + "." + fileExtension

		if fType != "image" {
			log.Println("Error not supported file:", fType+"/"+fSubType)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": "not supported: " + fType + "/" + fSubType})
			http.Error(w, "not supported: "+fType+"/"+fSubType, 415)
			return
		}

		// Crate file on disk
		out, err := os.Create("./upload/" + fileId + "." + fileExtension)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer out.Close()

		// Copy uploaded file contents to the file
		_, err = io.Copy(out, file)
		if err != nil {
			log.Println(err.Error())
			return
		}
		urls = append(urls, "/img/"+filename)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"status": "ok", "urls": urls})
}

func viewImageHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	file := c.URLParams["file"]
	http.ServeFile(w, r, "./upload/"+file)
}
