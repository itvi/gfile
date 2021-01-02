package main

import (
	"fmt"
	"gfile/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

// var root = "C:/test"

func main() {
	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", index)

	fmt.Println("start...")
	http.ListenAndServe(":9000", nil)
}

// file list
func index(w http.ResponseWriter, r *http.Request) {
	// name := r.URL.Query().Get("name")
	//isdir := r.URL.Query().Get("isdir")
	method := r.URL.Query().Get("method")
	path := r.URL.Query().Get("path") // http://localhost:9000/a

	funcMap := template.FuncMap{
		"cap": util.ConvertByteTo,
	}

	if method == "" {
		root := "C:/test"
		files := GetFiles(root)
		render(w, r, "./templates/index.html", funcMap, files)
	} else {
		files := GetFiles("C:/test" + path)
		fileName := "_list.html"
		filePath := "./templates/_list.html"
		renderPartial(w, fileName, filePath, funcMap, files)
	}
}

func GetFiles(dir string) []File {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	var myFiles []File
	for _, file := range files {
		path := strings.Replace(dir, "C:/test", "", -1) // TODO: const root
		f := File{
			Name:         file.Name(),
			IsDir:        file.IsDir(),
			Size:         file.Size(),
			LastModified: file.ModTime(),
			Path:         path,
		}
		myFiles = append(myFiles, f)
	}
	return myFiles
}
