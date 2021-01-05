package util

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// ConvertByteToMB ...
func ConvertByteTo(n int64) string {
	switch {
	case n == 0:
		return ""
	case n < 1024*1024:
		return fmt.Sprintf("%.2f", float64(n)/1024) + "KB"
	case n < 1024*1024*1024:
		return fmt.Sprintf("%.2f", float64(n)/1024/1024) + "MB"
	default:
		return fmt.Sprintf("%.2f", float64(n)/1024/1024/1024) + "GB"
	}
}

func Render(w http.ResponseWriter, r *http.Request, name string,
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

func RenderPartial(w http.ResponseWriter, fileName, filePath string,
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
