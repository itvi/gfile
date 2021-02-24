package handler

import (
	"gfile/internal/model"
	e "gfile/pkg/error"
	"gfile/pkg/form"
	"log"
	"net/http"
	"strconv"
)

type RoleHandler struct {
	M *model.RoleModel
}

func (rl *RoleHandler) Index(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := rl.M.GetRoles()
		if err != nil {
			log.Println("Get roles error:", err)
			return
		}
		otherTemplates := []string{
			p + "web/template/partial/toolbar_crud.html",
		}
		data := &TemplateData{
			Roles: roles,
		}
		c.render(w, r, otherTemplates, p+"web/template/html/role/index.html", data)
	}
}

func (rl *RoleHandler) AddView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := form.New(nil)
		data := &TemplateData{Form: form}
		c.render(w, r, nil, p+"web/template/html/role/add.html", data)
	}
}

func (rl *RoleHandler) Add(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			return
		}

		// if the form isn't valid, redisplay the template
		form := form.New(r.PostForm)
		form.Required("name")
		form.MaxLength("name", 20)
		form.MaxLength("desc", 50)
		data := &TemplateData{Form: form}
		tplName := p + "web/template/html/role/add.html"

		if !form.Valid() {
			c.render(w, r, nil, tplName, data)
			return
		}

		// create role
		name := form.Get("name")
		desc := form.Get("desc")
		role := &model.Role{Name: name, Description: desc}

		err := rl.M.Create(role)
		if err == e.ErrDuplicate {
			form.Errors.Add("name", "角色已存在")
			c.render(w, r, nil, tplName, data)
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		http.Redirect(w, r, "/roles", http.StatusSeeOther)
	}
}

func (rl *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println("Role ID convert error:", err)
		return
	}
	if err = rl.M.Delete(id); err != nil {
		log.Println(err)
		return
	}
}

func (rl *RoleHandler) EditView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil {
			log.Println("Role ID convert error:", err)
			return
		}

		role, err := rl.M.GetRole(id)
		if err == e.ErrNoRecord {
			log.Println("No record found")
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		form := form.New(nil)
		data := &TemplateData{
			Form: form,
			Role: role,
		}

		c.render(w, r, nil, p+"web/template/html/role/edit.html", data)
	}
}

func (rl *RoleHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println("Role ID convert error:", err)
		return
	}

	name := r.PostFormValue("name")
	desc := r.PostFormValue("desc")
	role := &model.Role{ID: id, Name: name, Description: desc}

	if err = rl.M.Edit(role); err != nil {
		log.Println(err)
		return
	}
}
