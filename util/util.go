package util

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// https://stackoverflow.com/questions/49057032/recursively-zipping-a-directory-in-golang
func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
