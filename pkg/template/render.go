package template

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplates(w http.ResponseWriter, r *http.Request,
	layouts []string, templateName, layoutName string,
	funcMaps map[string]interface{}, d interface{}) {
	layouts = append(layouts, templateName)

	// parse files
	t, err := template.New(templateName).Funcs(funcMaps).ParseFiles(layouts...)
	if err != nil {
		log.Println("parse files error:", err)
		w.Write([]byte(err.Error()))
	}

	// execute template
	err = t.ExecuteTemplate(w, layoutName, d)
	if err != nil {
		log.Println("execute error:", err)
		w.Write([]byte(err.Error()))
	}
}
