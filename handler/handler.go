package handler

import (
	"fmt"
	"gfile/model"
	"gfile/util"
	"net/http"
	"strconv"
	"text/template"
)

// file list
func Index(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// name := r.URL.Query().Get("name")
		//isdir := r.URL.Query().Get("isdir")
		method := r.URL.Query().Get("method")
		path := r.URL.Query().Get("path")

		funcMap := template.FuncMap{
			"cap": util.ConvertByteTo,
		}

		if method == "" {
			files := model.GetFiles(dir, dir)
			util.Render(w, r, "./templates/index.html", funcMap, files)
		} else {
			files := model.GetFiles(dir, dir+path)
			fileName := "_list.html"
			filePath := "./templates/_list.html"
			util.RenderPartial(w, fileName, filePath, funcMap, files)
		}
	}
}

func Download(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		path := r.URL.Query().Get("path")
		isdir := r.URL.Query().Get("isdir")

		var file string
		if isdir == "true" {
			file = name
		} else {
			file = dir + path + "/" + name
		}
		fmt.Println("Download file:", file)

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(name))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, file)
	}
}

func Zip(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	path := r.URL.Query().Get("path")
	//fmt.Println("path:", path)

	zipName := "./zip/" + name + ".zip"

	err := util.RecursiveZip(path, zipName)
	if err != nil {
		fmt.Println("zip error:", err)
		return
	}

	w.Write([]byte(zipName))
}
