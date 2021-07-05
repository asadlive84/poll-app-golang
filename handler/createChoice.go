package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"poll-app/storage"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type (
	ChoiceTempData struct {
		QuesionDetails storage.Question
		CSRFField      template.HTML
		Form           storage.Choice
		FormErrors     map[string]string
	}
)

func (s *Server) getCreateChoicePage(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	questionId, err := strconv.Atoi(id)

	if err != nil {
		log.Println(err)
	}

	question, err := s.store.GetQuestionDetail(int32(questionId))

	if err != nil {
		log.Fatalln(err)
	}

	formData := ChoiceTempData{
		QuesionDetails: *question,
		CSRFField:      csrf.TemplateField(r),
	}

	s.ChoiceCreateTemplate(w, r, formData)

}

func (s *Server) postCreateChoicePage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	questionId, err := strconv.Atoi(id)

	if err != nil {
		log.Println(err)
	}

	question, err := s.store.GetQuestionDetail(int32(questionId))

	if err != nil {
		log.Fatalln(err)
	}
	if err := r.ParseForm(); err != nil {
		log.Fatalln(err)
	}

	var form storage.Choice
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decode Error")
	}

	if err := form.ValidateChoice(); err != nil {
		vErros := map[string]string{}

		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErros[key] = value.Error()
				}
			}
		}
		data := ChoiceTempData{
			QuesionDetails: *question,
			CSRFField:      csrf.TemplateField(r),
			Form:           form,
			FormErrors:     vErros,
		}

		s.ChoiceCreateTemplate(w, r, data)
		return
	}

	session, _ := s.session.Get(r, "poll_app")
	user_email := session.Values["user_email"]

	user_id, err := s.store.GetUserDBID(user_email)

	if err != nil {
		log.Fatalln(err)
	}

	form.UserID = user_id
	form.QuestionID = question.ID
	log.Println("===========================choice==============")
	log.Println(form)

	a, b := s.store.SaveChoiceDB(form)
	//a refer for id
	//b refer fo err
	if b != nil {
		log.Println("data not saved")
	}

	log.Println(a)

	url := fmt.Sprintf("/get/question/%d/", question.ID)

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (s *Server) ChoiceCreateTemplate(w http.ResponseWriter, r *http.Request, form ChoiceTempData) {
	temp := s.templates.Lookup("create-choice.html")

	if err := temp.Execute(w, form); err != nil {
		log.Fatalln("executing template: ", err)
		return
	}
}
