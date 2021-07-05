package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"poll-app/storage"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/csrf"
)

type (
	Login struct {
		Email    string
		Password string
	}

	LoginTempData struct {
		CSRFField  template.HTML
		Form       storage.User
		FormErrors map[string]string
	}
)

func (l Login) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email,
			validation.Required.Error("Email dete hobe"),
			is.Email,
		),
		validation.Field(&l.Password,
			validation.Required.Error("Password is dete hobe"),
			validation.Length(3, 10).Error("Password Lenght must be 3 to 10"),
		),
	)
}

func (s *Server) getLoginPage(w http.ResponseWriter, r *http.Request) {

	formData := LoginTempData{
		CSRFField: csrf.TemplateField(r),
	}

	session, _ := s.session.Get(r, "poll_app")
	userId:=session.Values["user_id"] 

	if _,ok:=userId.(string);ok{
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}


	s.LoginTemplate(w, r, formData)

}

func (s *Server) postLoginPage(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Fatalln(err)
	}

	var form storage.User
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decode Error")
	}

	//log.Println("====================")

	if err := form.ValidateUser(); err != nil {
		vErros := map[string]string{}
		log.Println("validate error")
		log.Println(err)
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErros[key] = value.Error()
				}
			}
		}
		data := LoginTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErros,
		}

		s.LoginTemplate(w, r, data)
		return
	}

	fmt.Println(form)
	log.Println("=====================hello=============")
	log.Println(form)

	user, err := s.store.GetUser(form.Email, form.Password)

	if err != nil {
		log.Fatalln("user not found")
		return
	}

	log.Println(user)

	session, _ := s.session.Get(r, "poll_app")
	session.Values["user_id"] = strconv.Itoa(int(user.ID))
	session.Values["user_email"] = user.Email
	if err := session.Save(r, w); err != nil {
		log.Fatalln("saving error session")
	}

	// log.Println("=======session======")
	// log.Println(session)
	// get_seesion := session.Values["user_id"]
	// log.Println(get_seesion)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (s *Server) LoginTemplate(w http.ResponseWriter, r *http.Request, form LoginTempData) {
	temp := s.templates.Lookup("login.html")

	if err := temp.Execute(w, form); err != nil {
		log.Fatalln("executing template: ", err)
		return
	}
}
