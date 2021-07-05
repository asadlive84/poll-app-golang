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
	ChoiceUpdateTempData struct {
		QuesionDetails storage.Question
		Choice         storage.Choice
		CSRFField      template.HTML
		Form           storage.Choice
		FormErrors     map[string]string
	}
)

func (s *Server) getUpdateChoicePage(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	cId := mux.Vars(r)["cId"]

	questionId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}
	choiceId, err := strconv.Atoi(cId)

	if err != nil {
		log.Println(err)
	}

	question, err := s.store.GetQuestionDetail(int32(questionId))
	if err != nil {
		log.Fatalln(err)
	}
	choice, err := s.store.GetChoiceDetail(int32(choiceId), int32(questionId))

	if err != nil {
		log.Fatalln(err)
	}

	formData := ChoiceUpdateTempData{
		QuesionDetails: *question,
		CSRFField:      csrf.TemplateField(r),
		Choice:         *choice,
	}

	s.ChoiceUpdateTemplate(w, r, formData)

}

func (s *Server) postUpdateChoicePage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	cId := mux.Vars(r)["cId"]

	questionId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}
	choiceId, err := strconv.Atoi(cId)

	if err != nil {
		log.Println(err)
	}

	question, err := s.store.GetQuestionDetail(int32(questionId))
	if err != nil {
		log.Fatalln(err)
	}
	choiceDetail, err := s.store.GetChoiceDetail(int32(choiceId), int32(questionId))

	if err != nil {
		log.Fatalln(err)
	}

	var form storage.Choice
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decode Error")
	}

	log.Println("==========Encode form=========")
	log.Println(form.ChoiceText)
	log.Println(form.ID)
	log.Println("==========end===========")

	if err := form.ValidateChoice(); err != nil {
		vErros := map[string]string{}

		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErros[key] = value.Error()
				}
			}
		}
		data := ChoiceUpdateTempData{
			QuesionDetails: *question,
			CSRFField:      csrf.TemplateField(r),
			Form:           form,
			FormErrors:     vErros,
			Choice:         *choiceDetail,
		}

		s.ChoiceUpdateTemplate(w, r, data)
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

	a, b := s.store.UpdateChoiceDB(form.ChoiceText, choiceDetail.ID)
	//a refer for id
	//b refer fo err
	if b != nil {
		log.Println("data not saved")
	}

	log.Println(a)

	url := fmt.Sprintf("/get/question/%d/", question.ID)

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (s *Server) ChoiceUpdateTemplate(w http.ResponseWriter, r *http.Request, form ChoiceUpdateTempData) {
	temp := s.templates.Lookup("choiceUpdate.html")

	if err := temp.Execute(w, form); err != nil {
		log.Fatalln("executing template: ", err)
		return
	}
}
