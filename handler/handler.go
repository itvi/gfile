package handler

import (
	"gfile/model"
	"gfile/util"
	"net/http"
	"text/template"
)

func A(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(dir))
	}
}

// file list
func Index(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// name := r.URL.Query().Get("name")
		//isdir := r.URL.Query().Get("isdir")
		method := r.URL.Query().Get("method")
		path := r.URL.Query().Get("path") // http://localhost:9000/a

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
