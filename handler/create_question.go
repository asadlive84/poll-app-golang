package handler

import (
	"html/template"
	"log"
	"net/http"
	"poll-app/storage"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
)

type (
	QuestionTempData struct {
		CSRFField  template.HTML
		Form       storage.Question
		FormErrors map[string]string
	}
)

func (s *Server) getCreateQuesionPage(w http.ResponseWriter, r *http.Request) {
	formData := QuestionTempData{
		CSRFField: csrf.TemplateField(r),
	}

	// session, _ := s.session.Get(r, "poll_app")
	// value := session.Values["user_id"]
	// log.Println("GEt the usr value==============")
	// log.Println(value)

	s.QuestionTemplate(w, r, formData)
}

func (s *Server) postCreateQuesionPage(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Fatalln(err)
	}

	var form storage.Question
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decode Error")
	}

	if err := form.ValidateQuestion(); err != nil {
		vErros := map[string]string{}

		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErros[key] = value.Error()
				}
			}
		}
		data := QuestionTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErros,
		}

		s.QuestionTemplate(w, r, data)
		return
	}

	session, _ := s.session.Get(r, "poll_app")
	user_email := session.Values["user_email"]

	user_id, err := s.store.GetUserDBID(user_email)

	if err != nil {
		log.Fatalln(err)
	}

	form.UserID = user_id
	id, err := s.store.SaveQuestionDB(form)
	if err != nil {
		log.Println("data not saved")
	}

	log.Println(id)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (s *Server) QuestionTemplate(w http.ResponseWriter, r *http.Request, form QuestionTempData) {
	temp := s.templates.Lookup("create-question.html")

	if err := temp.Execute(w, form); err != nil {
		log.Fatalln("executing template: ", err)
		return
	}
}
