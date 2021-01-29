package handler

import (
	"gfile/internal/model"
	"gfile/pkg/form"
	"log"
	"net/http"
	"strconv"
)

type CasbinHandler struct {
	M *model.CasbinModel
}

func (a *CasbinHandler) Index(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleName := r.URL.Query().Get(":name")
		policies := a.M.GetPoliciesOrderBy(roleName)
		otherTemplates := []string{
			"./web/template/partial/toolbar_crud.html",
		}
		data := &TemplateData{Policies: policies}
		c.render(w, r, otherTemplates, "./web/template/html/casbin/index.html", data)
	}
}

func (a *CasbinHandler) AddView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := form.New(nil)
		data := &TemplateData{Form: form}
		c.render(w, r, nil, "./web/template/html/casbin/add.html", data)
	}
}

func (a *CasbinHandler) Add(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	sub := r.PostFormValue("sub")
	obj := r.PostFormValue("obj")
	act := r.PostFormValue("act")

	// ref := r.Header["Referer"]
	// log.Println("ref:", ref)

	enforcer, err := a.M.Init()
	if err != nil {
		log.Fatal("init casbin error")
	}

	_, err = enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		log.Println("Add policy error:", err)
	}

	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

func (a *CasbinHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sub := r.PostFormValue("sub")
	obj := r.PostFormValue("obj")
	act := r.PostFormValue("act")

	enforcer, err := a.M.Init()
	if err != nil {
		log.Fatal("init casbin error")
	}
	_, err = enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		log.Println(err)
	}
}

func (a *CasbinHandler) EditView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub := r.URL.Query().Get("sub")
		obj := r.URL.Query().Get("obj")
		act := r.URL.Query().Get("act")
		log.Println(sub, obj, act)
		policy := &model.CasbinPolicy{
			Sub: sub,
			Obj: obj,
			Act: act,
		}
		form := form.New(nil)
		data := &TemplateData{
			Policy: policy,
			Form:   form,
		}

		c.render(w, r, nil, "./web/template/html/casbin/edit.html", data)
	}
}

func (a *CasbinHandler) Edit(w http.ResponseWriter, r *http.Request) {
	old := []string{
		r.PostFormValue("oSub"),
		r.PostFormValue("oObj"),
		r.PostFormValue("oAct"),
	}
	new := []string{
		r.PostFormValue("sub"),
		r.PostFormValue("obj"),
		r.PostFormValue("act"),
	}

	if err := a.M.Edit(old, new); err != nil {
		log.Println("update fault")
		return
	}
}

// AddRolesForUserView add roles for user (GET)
func (a *CasbinHandler) AddRolesForUserView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get(":userid"))
		if err != nil {
			log.Println("Convert userID error:", err)
			return
		}

		// get user by id
		userModel := &model.UserModel{DB: a.M.DB}
		user, err := userModel.GetUser(userID)
		if err != nil {
			log.Println(err)
			return
		}

		// get all roles
		roleModel := &model.RoleModel{DB: a.M.DB}
		roles, err := roleModel.GetRoles()
		if err != nil {
			log.Println(err)
			return
		}

		// get roles for user
		enforcer, err := a.M.Init()
		if err != nil {
			log.Fatal("init casbin error")
		}

		userRoles, err := enforcer.GetRolesForUser(user.SN)
		if err != nil {
			log.Println("Get roles error:", err)
			return
		}

		data := &TemplateData{
			User:         user,
			Roles:        roles,
			RolesForUser: userRoles,
		}

		c.render(w, r, nil, "./web/template/html/casbin/add_roles_for_user.html", data)
	}
}

// AddRolesForUser add roles for user (POST)
func (a *CasbinHandler) AddRolesForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get(":userid"))
	if err != nil {
		log.Println("Convert userID error:", err)
		return
	}

	if err = r.ParseForm(); err != nil {
		log.Fatal("Parse Form error,not found roles", err)
	}
	roles := r.Form["roles"]

	// get user by id
	userModel := &model.UserModel{DB: a.M.DB}
	user, err := userModel.GetUser(userID)
	if err != nil {
		log.Println(err)
		return
	}

	// add roles for user(1. delete; 2. add)
	enforcer, err := a.M.Init()
	if err != nil {
		log.Fatal("init casbin error")
	}

	_, err = enforcer.DeleteRolesForUser(user.SN)
	if err != nil {
		log.Fatal("Delete roles for user error:", err)
	}

	_, err = enforcer.AddRolesForUser(user.SN, roles)
	if err != nil {
		log.Fatal("Add roles for user error:", err)
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
