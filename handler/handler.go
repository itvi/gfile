package handler

import (
	"database/sql"
	"fmt"
	"gfile/model"
	"gfile/util"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/karrick/godirwalk"
)

// file list
func Index(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			filePath := "./templates/partial/_list.html"
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
			file = dir + path
		}
		fmt.Println("Download file:", file)

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(name))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, file)
	}
}

func Zip(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		path := r.URL.Query().Get("path")

		pathToZip := dir + path
		zipName := "./zip/" + name + ".zip"

		err := util.RecursiveZip(pathToZip, zipName)
		if err != nil {
			fmt.Println("zip error:", err)
			return
		}

		w.Write([]byte(zipName))
	}
}

func Search(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		s := `SELECT name,isdir,size,last_modified,path 
			  FROM files WHERE name LIKE '%` + q + `%'`
		rows, err := db.Query(s)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		var files []*model.File
		layout := "2006-01-02 15:04:05"

		for rows.Next() {
			file := &model.File{}
			var lastModified string
			if err = rows.Scan(&file.Name, &file.IsDir, &file.Size,
				&lastModified, &file.Path); err != nil {
				fmt.Println("rows scan error:", err)
				return
			}
			// parse datatime -> 2021-01-04 08:25:32.629566+08:00
			last := strings.Split(lastModified, "+")[0]
			lastModifiedDate, err := time.Parse(layout, last)
			if err != nil {
				fmt.Println(err)
				return
			}
			file.LastModified = lastModifiedDate
			files = append(files, file)
		}
		if err = rows.Err(); err != nil {
			return
		}
		funcMap := template.FuncMap{
			"cap": util.ConvertByteTo,
		}
		util.Render(w, r, "./templates/index.html", funcMap, files)
	}
}

func Rebuild(dir string, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// clear first
		stmt, err := db.Prepare("DELETE FROM files;")
		if err != nil {
			log.Println("delete prepare error:", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			log.Println("delete exec error:", err)
			return
		}

		rebuild(db, dir)
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
				pathName = "\\" + osPathname
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
