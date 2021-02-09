package handler

import (
	"context"
	"database/sql"
	"fmt"
	"gfile/internal/model"
	e "gfile/pkg/error"
	"gfile/pkg/form"
	tmp "gfile/pkg/template"
	"gfile/pkg/util"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type TemplateData struct {
	Now               time.Time
	AuthenticatedUser *model.User
	Form              *form.Form
	User              *model.User
	Users             []*model.User
	Role              *model.Role
	Roles             []*model.Role
	RolesForUser      []string
	Policy            *model.CasbinPolicy
	Policies          []*model.CasbinPolicy
	File              *model.File
	Files             []*model.File
	FileStat          map[string]int
}

type Configuration struct {
	Session *sessions.Session
	Home    *HomeHandler
	User    *UserHandler
	Role    *RoleHandler
	Casbin  *CasbinHandler
	File    *FileHandler
}

type contextKey string

var contextKeyUser = contextKey("user")

func Config(dir string) (*Configuration, *sql.DB) {
	// database
	db, err := openDB("./db.db")
	if err != nil {
		err = fmt.Errorf("open db error: %w -> from open db", err)
		log.Panic(err)
	}

	// session
	session := sessions.New([]byte("afkkjfkajf!23234324#@#$"))
	session.Lifetime = 15 * time.Minute

	c := &Configuration{
		Session: session,
		Home:    &HomeHandler{},
		User:    &UserHandler{M: &model.UserModel{DB: db}},
		Role:    &RoleHandler{M: &model.RoleModel{DB: db}},
		Casbin:  &CasbinHandler{M: &model.CasbinModel{DB: db}},
		File:    &FileHandler{M: &model.FileModel{DB: db}, Dir: dir},
	}
	return c, db
}

func (c *Configuration) authenticatedUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(contextKeyUser).(*model.User)
	if !ok {
		return nil
	}
	return user
}

func (c *Configuration) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/users/login", 302)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *Configuration) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if a "userID" exists in the session.
		exists := c.Session.Exists(r, "userID")
		if !exists {
			log.Println("session not exist.")
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the details of the current user from the database.
		userID := c.Session.GetInt(r, "userID")
		fmt.Println("session userid:", userID)
		user, err := c.User.M.GetUser(userID)
		if err == e.ErrNoRecord {
			c.Session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		// The request is coming form a valid, authenticated user.
		// Create a new copy of the request with the user
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Configuration) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enforcer, err := c.Casbin.M.Init()
		if err != nil {
			log.Println("casbin init error")
		}

		// get current logined user
		user := c.authenticatedUser(r)

		var sub string
		if user == nil {
			sub = "anonymous"
		} else {
			sub = user.SN
		}
		obj := r.URL.Path
		act := r.Method

		ok, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			log.Println("casbin enforce error:", err)
		}
		if ok {
			next.ServeHTTP(w, r)
		} else {
			log.Println("forbidden [", sub, obj, act+" ]")
			w.Write([]byte("forbidden"))
		}
	})
}

// database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Configuration) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}
	td.AuthenticatedUser = c.authenticatedUser(r)

	return td
}

// Render render templates with preset layouts
func (c *Configuration) render(w http.ResponseWriter, r *http.Request,
	otherTemplates []string, templateName string, data *TemplateData) {

	layouts := []string{
		"./web/template/layout.html",
		"./web/template/partial/menu.html",
	}
	layouts = append(layouts, otherTemplates...)

	templateData := c.addDefaultData(data, r)
	tmp.RenderTemplates(w, r, layouts, templateName, "layout", funcMaps, templateData)
}

var funcMaps = template.FuncMap{
	"safe": func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"ownedRoles": func(role string, roles []string) string {
		for _, r := range roles {
			if r == role {
				return "checked"
			}
		}
		return ""
	},
	"cap": util.ConvertByteTo,
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

func Render(w http.ResponseWriter, r *http.Request, name string,
	funcMaps map[string]interface{}, d interface{}) {
	baseFile := "layout"
	tmpls := []string{
		"./web/template/layout.html",
		"./web/template/partial/menu.html",
		"./web/template/partial/breadcrumb.html",
		"./web/template/partial/toolbar.html",
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

// // watcher(I want to batch delete/add to database. How to pass channel's item to a slice?)
// func monitor(dir string) {
// 	var files []string

// 	createFiles := []model.File{}

// 	go func() {
// 		w := watcher.New()
// 		w.FilterOps(watcher.Create, watcher.Remove, watcher.Rename, watcher.Move)

// 		go func() {
// 			var i int
// 			for {
// 				select {
// 				case event := <-w.Event:
// 					//fmt.Println(event)

// 					removeChan := make(chan string)
// 					createChan := make(chan model.File)
// 					switch event.Op {
// 					case watcher.Remove:
// 						file := event.Path[len(dir):]
// 						//fmt.Println("File:", file)

// 						go func(file string) {
// 							removeChan <- file
// 							close(removeChan)
// 						}(file)

// 						for n := range removeChan {
// 							files = append(files, n)
// 						}
// 						// delete from database:
// 						fmt.Printf("del slice count:%d, %v\n", len(files), files)
// 						fmt.Println("delete from database")
// 						//config.File.M.DeleteIndexes(files)
// 					case watcher.Create:
// 						i++
// 						fmt.Println("Create.", i)

// 						file := &model.File{
// 							Name: event.Name(),
// 							Size: event.Size(),
// 							Path: event.Path[len(dir):],
// 						}

// 						go func(file model.File) {
// 							createChan <- file
// 							close(createChan)
// 						}(*file)

// 						for f := range createChan {
// 							createFiles = append(createFiles, f)
// 						}
// 						fmt.Println("create:", createFiles)
// 						// add to database

// 					default:
// 						fmt.Println("default")
// 					}
// 					// do something from database

// 				case err := <-w.Error:
// 					log.Println(err)
// 				case <-w.Closed:
// 					return
// 				}
// 			}
// 		}()
// 		fmt.Print("O")

// 		if err := w.AddRecursive(dir); err != nil {
// 			log.Fatalln(err)
// 		}
// 		if err := w.Start(time.Millisecond * 100); err != nil {
// 			log.Fatalln(err)
// 		}
// 	}()
// }

// func fromChannel(file <-chan string) []string {
// 	var results []string
// 	for {
// 		select {
// 		case str := <-file:
// 			results = append(results, str)
// 			if len(results) == 3 {
// 				fmt.Println("Got all results")
// 				return results
// 			}
// 		}
// 	}
// }

// func putInChannel(file string, channel chan<- string) {
// 	channel <- file
// }

// copy finished
// https://github.com/radovskyb/watcher/issues/46
