package handler

import (
	"gfile/internal/model"
	e "gfile/pkg/error"
	"gfile/pkg/forms"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	M *model.UserModel
}

func (u *UserHandler) Index(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := u.M.GetUsers()
		if err != nil {
			log.Println("Get users error:", err)
			return
		}
		otherTemplates := []string{
			"./web/template/partial/toolbar_crud.html",
		}
		data := &TemplateData{
			Users: users,
		}
		c.render(w, r, otherTemplates, "./web/template/html/user/index.html", data)
	}
}

func (u *UserHandler) AddView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := forms.New(nil)
		data := &TemplateData{Form: form}
		c.render(w, r, nil, "./web/template/html/user/add.html", data)
	}
}

func (u *UserHandler) Add(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}

		// validate
		form := forms.New(r.PostForm)
		form.Required("sn", "password")
		data := &TemplateData{Form: form}
		page := "./web/template/html/user/add.html"

		if !form.Valid() {
			c.render(w, r, nil, page, data)
			return
		}

		// create user
		var user = &model.User{
			SN:       form.Get("sn"),
			Name:     form.Get("name"),
			Email:    form.Get("email"),
			Password: form.Get("password"),
		}

		err = u.M.Create(user)
		if err == e.ErrDuplicate {
			form.Errors.Add("sn", "用户已存在")
			c.render(w, r, nil, page, data)
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
}

func (u *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println("User ID convert error:", err)
		return
	}

	err = u.M.Delete(id)
	if err != nil {
		log.Println(err)
		return
	}
}

func (u *UserHandler) EditView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil {
			log.Println("User ID convert error:", err)
			return
		}

		// TODO: maybe in middleware
		// don't change others data
		// userID := c.Session.GetInt(r, "userID")
		// if id != userID {
		// 	w.Write([]byte("forbidden"))
		// 	return
		// }

		user, err := u.M.GetUser(id)
		if err == e.ErrNoRecord {
			log.Println("No record found")
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		form := forms.New(nil)
		data := &TemplateData{
			Form: form,
			User: user,
		}
		c.render(w, r, nil, "./web/template/html/user/edit.html", data)
	}
}

func (u *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println("User ID convert error:", err)
		return
	}

	sn := r.PostFormValue("sn")
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	user := &model.User{
		ID:             id,
		SN:             sn,
		Name:           name,
		Email:          email,
		HashedPassword: []byte(password),
	}

	if err = u.M.Edit(user); err != nil {
		log.Println(err)
		return
	}
}

func (u *UserHandler) LoginView(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := forms.New(nil)
		data := &TemplateData{Form: form}
		c.render(w, r, nil, "./web/template/html/user/login.html", data)
	}
}

// Login use configuration as parameter is for set session
func (u *UserHandler) Login(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println("Parse form error:", err)
			return
		}

		form := forms.New(r.PostForm)
		sn := form.Get("sn")
		password := form.Get("password")

		user, err := c.User.M.Authenticate(sn, password)
		if err == e.ErrInvalidCredentials {
			form.Errors.Add("generic", "用户或密码不正确！")
			data := &TemplateData{
				Form: form,
			}
			c.render(w, r, nil, "./web/template/html/user/login.html", data)
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		// session
		c.Session.Put(r, "userID", user.ID)

		// // create a new jwt token
		// tokenString, err := auth.SignToken(string(user.SN))
		// if err != nil {
		// 	log.Printf("User %s generate token error: %v", user.Name, err)
		// 	return
		// }
		// // create cookie
		// exp := time.Now().Add(15 * time.Minute)

		// cookie := &http.Cookie{
		// 	Name:     "token",
		// 	Value:    tokenString,
		// 	Expires:  exp,
		// 	HttpOnly: true,
		// 	Secure:   true,
		// }
		// http.SetCookie(w, cookie)

		// // custom claims
		// type MyClaims struct {
		// 	Name string
		// 	jwt.StandardClaims
		// }

		// claims := MyClaims{
		// 	user.SN,
		// 	jwt.StandardClaims{
		// 		ExpiresAt: 15000,
		// 		Issuer:    "testtt",
		// 	},
		// }
		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// s, err := token.SignedString(auth.JWTTOKEN)
		// fmt.Printf("%v %v", s, err)

		// http.SetCookie(w, &http.Cookie{
		// 	Name:     "token123",
		// 	Value:    s,
		// 	Domain:   "",
		// 	HttpOnly: true,
		// })

		// =========================================
		//	r.Header.Set("Authorization", s)
		// fmt.Println(tokenString)
		// httputil.DumpResponse(w, tokenString)

		//w.Header().Set("Authorization", tokenString)
		// r.Header.Add("Authorization", tokenString)
		// r.Header.Set("Authorization", tokenString)
		// fmt.Println("login token:", tokenString)
		//r.Header.Set("Authorization", tokenString) //
		//r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenString))

		//w.Write([]byte(r.Header.Get("Autorrization")))
		//w.Write([]byte(tokenString)) //

		// var bearer = tokenString

		// req, err := http.NewRequest("GET", "http://localhost:9000/roles", nil)
		// req.Header.Add("Authorization", bearer)

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	log.Println("client err:", err)
		// }
		// body, _ := ioutil.ReadAll(resp.Body)
		// log.Println(string([]byte(body)))

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout remove session
func (u *UserHandler) Logout(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Session.Remove(r, "userID")
		http.Redirect(w, r, "/", 303)
	}
}
