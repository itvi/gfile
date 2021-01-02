package main

import (
	"fmt"
	"gfile/util"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

type File struct {
	Name         string
	IsDir        bool
	Size         int64
	LastModified time.Time
	Path         string
}

func render(w http.ResponseWriter, r *http.Request, name string,
	funcMaps map[string]interface{}, d interface{}) {
	baseFile := "layout"
	tmpls := []string{
		"./templates/layout.html",
		"./templates/menu.html",
	}
	tmpls = append(tmpls, name)

	// parse files
	t, err := template.New(name).Funcs(funcMaps).ParseFiles(tmpls...)
	if err != nil {
		log.Println("parse files error:", err)
		w.Write([]byte(err.Error()))
	}

	// execute template
	err = t.ExecuteTemplate(w, baseFile, d)
	if err != nil {
		log.Println("execute error:", err)
		w.Write([]byte(err.Error()))
	}
}

func renderPartial(w http.ResponseWriter, fileName, filePath string,
	funcMap map[string]interface{}, data interface{}) {
	t, err := template.New(fileName).Funcs(funcMap).ParseFiles(filePath)
	if err != nil {
		log.Println("Parse file error:", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println("Execute template error:", err)
	}
}

var root = "C:/test"

func main() {
	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/dl", download)
	http.HandleFunc("/open", open)

	fmt.Println("start...")
	http.ListenAndServe(":9000", nil)
}

// file list
func index(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	isdir := r.URL.Query().Get("isdir")
	method := r.URL.Query().Get("method")
	href := r.URL.Query().Get("href")

	funcMap := template.FuncMap{
		"cap": util.ConvertByteTo,
	}

	if method == "" {
		root = "C:/test"
		files := GetFiles(root)
		render(w, r, "./templates/index.html", funcMap, files)
	} else {
		if isdir == "dir" {
			root = root + "/" + name
		}
		fmt.Println("href:", href)
		u, _ := url.Parse(href)
		// root = root + u.Path + name
		fmt.Println("root:", u.Path)
		files := GetFiles(root)
		fileName := "_list.html"
		filePath := "./templates/_list.html"
		renderPartial(w, fileName, filePath, funcMap, files)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	files := GetFiles(root)
	render(w, r, "./templates/download.html", nil, files)
}

func GetFiles(dir string) []File {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	var myFiles []File
	for _, file := range files {
		f := File{
			Name:         file.Name(),
			IsDir:        file.IsDir(),
			Size:         file.Size(),
			LastModified: file.ModTime(),
			Path:         dir,
		}
		myFiles = append(myFiles, f)
	}
	return myFiles
}

func open(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("dir"))
	path := r.URL.Query().Get("path")
	files := GetFiles(root + "/" + path)
	funcMap := template.FuncMap{
		"cap": util.ConvertByteTo,
	}
	render(w, r, "./templates/index.html", funcMap, files)
}
