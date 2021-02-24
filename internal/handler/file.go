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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/radovskyb/watcher"
)

type FileHandler struct {
	Dir string
	M   *model.FileModel
}

// file list
func (f *FileHandler) Index(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		// file stat
		stat := f.M.FileStat(name)
		fmt.Println(stat)

		//if method == "" {
		files := model.GetFiles(f.Dir, f.Dir)
		otmps := []string{
			p + "web/template/partial/breadcrumb.html",
			p + "web/template/partial/toolbar.html",
		}
		c.render(w, r, otmps, p+"web/template/html/file/index.html", &TemplateData{
			Files:    files,
			FileStat: stat,
		})
	}
}

func (f *FileHandler) getDirContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		files := model.GetFiles(f.Dir, f.Dir+path)
		var dirNum, fileNum int
		for _, file := range files {
			if file.IsDir {
				dirNum++
			}
			if !file.IsDir {
				fileNum++
			}
		}
		var stat = make(map[string]int)
		stat["dir"] = dirNum
		stat["file"] = fileNum

		fileName := "list.html"
		filePath := p + "web/template/partial/list.html"

		funcMap := template.FuncMap{
			"cap": util.ConvertByteTo,
		}

		data := struct {
			File     []*model.File
			FileStat map[string]int
		}{File: files, FileStat: stat}

		fmt.Printf("%T", data)
		RenderPartial(w, fileName, filePath, funcMap, data)
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

		pathToZip := p + path
		zipName := p + "zip/" + name + ".zip"
		fmt.Printf("pathtozip: %s, zipname:%s\n", pathToZip, zipName)

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

		// file stat
		stat := f.M.FileStat(q)

		otmps := []string{
			p + "web/template/partial/breadcrumb.html",
			p + "web/template/partial/toolbar.html",
		}
		c.render(w, r, otmps, p+"web/template/html/file/search.html", &TemplateData{
			Files:    files,
			FileStat: stat,
		})
	}
}

func (f *FileHandler) Rebuild(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// clear first
		err := f.M.ClearFileIndexes()
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
				pathName = string(os.PathSeparator) + osPathname
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

func (f *FileHandler) Watchdog(dir string) {

	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Remove, watcher.Rename, watcher.Move)

	go func() {
		var i int
		for {
			select {
			case event := <-w.Event:
				//fmt.Println(event)

				switch event.Op {

				case watcher.Remove:
					file := event.Path[len(dir):]
					log.Println("Will remove:", file)
					// delete from database:
					if err := f.M.DeleteIndex(file); err != nil {
						log.Println("Delete index error:", err)
					}

				case watcher.Create:
					i++
					fmt.Println("Create.", i)

					file := &model.File{
						Name:  event.Name(),
						IsDir: event.IsDir(),
						Size:  event.Size(),
						Path:  event.Path[len(dir):],
					}
					//fmt.Printf("File will create: %v\n", file)
					// add to database
					if err := f.M.CreateIndex(file); err != nil {
						log.Println("Add index error:", err)
					}
				case watcher.Rename:
					oldName := filepath.Base(event.OldPath)
					newName := filepath.Base(event.Path)
					oldFile := &model.File{
						Name: oldName,
						Path: event.OldPath[len(dir):],
					}
					newFile := &model.File{
						Name: newName,
						Path: event.Path[len(dir):],
					}
					//fmt.Printf("old %v, new %v", oldFile, newFile)
					if err := f.M.UpdateIndex(oldFile, newFile); err != nil {
						log.Println("Rename index error:", err)
					}
				default:
					fmt.Println("default")
				}

			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()
	// fmt.Print("O")

	if err := w.AddRecursive(dir); err != nil {
		log.Fatalln(err)
	}
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
