package handler

import (
	"database/sql"
	"fmt"
	"gfile/internal/model"
	"gfile/pkg/util"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
)

type FileHandler struct {
	Dir string
	M   *model.FileModel
}

// file list
func (f *FileHandler) Index(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		path := r.URL.Query().Get("path")
		fmt.Println("method:", method)

		if method == "" {
			files := model.GetFiles(f.Dir, f.Dir)
			otmps := []string{
				"./web/template/partial/breadcrumb.html",
				"./web/template/partial/toolbar.html",
			}
			c.render(w, r, otmps, "./web/template/html/file/index.html", &TemplateData{
				Files: files,
			})
		} else {
			files := model.GetFiles(f.Dir, f.Dir+path)
			fileName := "list.html"
			filePath := "./web/template/partial/list.html"

			funcMap := template.FuncMap{
				"cap": util.ConvertByteTo,
			}
			RenderPartial(w, fileName, filePath, funcMap, files)
		}
	}
}

func (f *FileHandler) Download(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		path := r.URL.Query().Get("path")
		isdir := r.URL.Query().Get("isdir")

		var file string
		if isdir == "true" {
			file = name
		} else {
			file = f.Dir + path
		}
		fmt.Println("Download file:", file)

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(name))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, file)
	}
}

func (f *FileHandler) Zip(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		path := r.URL.Query().Get("path")

		pathToZip := f.Dir + path
		zipName := "./zip/" + name + ".zip"

		err := util.RecursiveZip(pathToZip, zipName)
		if err != nil {
			fmt.Println("zip error:", err)
			return
		}

		w.Write([]byte(zipName))
	}
}

func (f *FileHandler) Search(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")

		files, err := f.M.Search(q)
		if err != nil {
			log.Println("Get files error:", err)
			return
		}

		otmps := []string{
			"./web/template/partial/breadcrumb.html",
			"./web/template/partial/toolbar.html",
		}
		c.render(w, r, otmps, "./web/template/html/file/search.html", &TemplateData{
			Files: files,
		})
	}
}

func (f *FileHandler) Rebuild(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("rebild")
		// clear first
		err := f.M.DeleteFileIndex()
		if err != nil {
			log.Println("Clear index error:", err)
			return
		}

		rebuild(f.M.DB, f.Dir)
		w.Write([]byte("刷新成功！"))
	}
}

func rebuild(db *sql.DB, dir string) {
	start := time.Now()

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}

	count := 0
	var name string
	var isdir bool
	var lastModified time.Time

	godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			// skip specified folder
			if strings.Contains(osPathname, "zip") {
				return godirwalk.SkipThis
			}

			count++
			stat, err := os.Stat(osPathname)
			if err != nil {
				return err
			}
			name = de.Name()
			isdir = de.IsDir()

			size := stat.Size()
			lastModified = stat.ModTime()

			var pathName string
			// ./static/css
			if dir == "." { // static\css
				//dir = ""
				pathName = "/" + osPathname
			} else {
				pathName = osPathname[len(dir):]
			}
			// fmt.Println(count, pathName, name, isdir, size, lastModified)

			// add to database
			s := `INSERT INTO files(name,isdir,size,last_modified,path 
				)VALUES(?,?,?,?,?)`
			_, err = tx.Exec(s, name, isdir, size, lastModified, pathName)
			if err != nil {
				tx.Rollback()
				fmt.Println("insert error:", err)
			}

			return err
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})

	tx.Commit()
	fmt.Printf("Rebuild: %d items, Spend: %s\n", count, time.Since(start))
}
