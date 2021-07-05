package handler

import (
	"fmt"
	"log"
	"net/http"
	"poll-app/storage"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func (s *Server) getUpdateQuestionPage(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	questionId, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	question, err := s.store.GetQuestionDetail(int32(questionId))
	if err != nil {
		log.Println("query doesnot match")
	}

	formData := QuestionTempData{
		CSRFField: csrf.TemplateField(r),
		Form:      *question,
	}

	s.UpdateQuestionTemplate(w, r, formData)
}

func (s *Server) postUpdateQuestionPage(w http.ResponseWriter, r *http.Request) {

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
	id, err := s.store.UpdateQuestionDB(form.QuestionText, form.UserID, form.ID)
	if err != nil {
		log.Println("data not updated")
	}
	log.Println(id)
	

	redirect_url := fmt.Sprintf("/get/question/%d/", form.ID)



	http.Redirect(w, r, redirect_url, http.StatusSeeOther)

}

func (s *Server) UpdateQuestionTemplate(w http.ResponseWriter, r *http.Request, form QuestionTempData) {
	temp := s.templates.Lookup("update-question.html")

	if err := temp.Execute(w, form); err != nil {
		log.Fatalln("executing template: ", err)
		return
	}
}
