package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates := template.Must(template.ParseGlob("./templates/*.tmpl"))
		if err := templates.ExecuteTemplate(w, "login_page", nil); err != nil {
			log.Println("Error rendering template:", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if r.Method == "POST" {
		log.Println(r.FormValue("username"), r.FormValue("password"))
		http.Redirect(w, r, "/index/", http.StatusFound)
	}
}

func articleViewHandler(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["article"]
	var article Article
	err := db.Read("article", articleId, &article)
	if err != nil {
		http.Redirect(w, r, "/"+articleId+"/edit", http.StatusFound)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"status": "ok", "url": "/" + articleId})
		return
	}

	templates := template.Must(template.ParseGlob("./templates/*.tmpl"))
	err = templates.ExecuteTemplate(w, "view_page", article)
	if err != nil {
		log.Println("Error rendering template:", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

func articleEditHandler(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["article"]
	if r.Method == "GET" {
		var article Article
		err := db.Read("article", articleId, &article)
		if err != nil {
			article = Article{}
			article.Title = "new page"
			article.Content = ""
			article.Comments = false
			article.Markdown = ""
			if err := db.Write("article", articleId, &article); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
		templates := template.Must(template.ParseGlob("./templates/*.tmpl"))
		err = templates.ExecuteTemplate(w, "edit_page", article)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else if r.Method == "POST" {
		var article Article
		err := db.Read("article", articleId, &article)
		if err != nil {
			log.Println("Error reading body:", err.Error())
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": err.Error()})
			http.Error(w, err.Error(), 500)
			return
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

		article.Markdown = string(form["content"])
		unsafe := blackfriday.MarkdownCommon([]byte(article.Markdown))
		html := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
		article.Content = html

		if err := db.Write("article", articleId, &article); err != nil {
			log.Println("Error saving article:", err.Error())
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": err.Error()})
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"status": "ok", "url": "/" + articleId + "/"})
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

func uploadPostHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println(fType, fSubType, filename)

		if fType != "image" {
			log.Println("Error not supported file:", fType+"/"+fSubType)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{"status": "error", "msg": "not supported: " + fType + "/" + fSubType})
			http.Error(w, "not supported: "+fType+"/"+fSubType, 415)
			return
		}

		// Crate file on disk
		out, err := os.Create("./upload/" + filename)
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
		//urls = append(urls, "./img/"+filename)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"status": "ok", "urls": urls})
}

func viewImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	http.ServeFile(w, r, "./upload/"+vars["file"])
}
