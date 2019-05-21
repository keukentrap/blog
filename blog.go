package main

import (
	"fmt"
	"html/template"
	"log"
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gomarkdown/markdown"
)

var tpl struct {
	index, article *template.Template
}
var docs map[string]template.HTML

type Context struct {
	Article template.HTML
	Request *http.Request
}


func init() {

	TemplatePath := "templates/"
	base := filepath.Join(TemplatePath, "base.gohtml")
	index := filepath.Join(TemplatePath, "index.gohtml")
	article := filepath.Join(TemplatePath, "article.gohtml")
	tpl.index = template.Must(template.ParseFiles(base, index))
	tpl.article = template.Must(template.ParseFiles(base, article))

	files, err := ioutil.ReadDir("articles/")
	if err != nil {
		log.Fatalln(err)
	}

	docs = make(map[string]template.HTML)
	for _, file := range files {
		md , err := ioutil.ReadFile("articles/" + file.Name())
		if err != nil {
			log.Fatalln(err)
		}
		n := file.Name()[:len(file.Name())-3]
		docs[n] = template.HTML(markdown.ToHTML(md, nil, nil))
	}
}

func main() {
	ip_port := ":8080"
	http.HandleFunc("/",handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Serving from " + ip_port)
	http.ListenAndServe(ip_port, nil)
}

func handler(conn http.ResponseWriter,req *http.Request) {

	var buf  bytes.Buffer
	log.Println(req.URL.String())
	if val , ok := docs[req.URL.String()[1:]]; ok {
		context := Context{val,req}
		err := tpl.article.ExecuteTemplate(&buf, "base", context)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := tpl.index.ExecuteTemplate(&buf, "base", docs)
		if err != nil {
			log.Fatalln(err)
		}
	}
	body := buf.String()

	fmt.Fprintf(conn, body)
}
