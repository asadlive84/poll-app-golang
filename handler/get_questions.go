package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"poll-app/storage"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type (
	TemplateData struct {
		Question []storage.Question
	}

	TempDataDetails struct {
		CSRFField template.HTML
		Question  storage.Question
		Choices   []storage.Choice
		UserId    int32
	}
)

func (s *Server) getAllQuestionPage(w http.ResponseWriter, r *http.Request) {

	questionDb, err := s.store.GetAllQuestionDB()
	if err != nil {
		log.Println("database error")
	}
	temp := TemplateData{
		Question: questionDb,
	}

	s.DefaultTemplate(w, r, "get_question.html", temp)
}

func (s *Server) questionDetails(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	questionId, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	question, err := s.store.GetQuestionDetail(int32(questionId))
	if err != nil {
		log.Println("query doesnot match")
	}

	choiceSet, err := s.store.ChoiceQuery(question.ID)

	if err != nil {
		log.Fatalln(err)
	}

	session, _ := s.session.Get(r, "poll_app")
	user_email := session.Values["user_email"]

	user_id, err := s.store.GetUserDBID(user_email)

	if err != nil {
		log.Fatalln(err)
	}

	temp := TempDataDetails{
		CSRFField: csrf.TemplateField(r),
		Question:  *question,
		Choices:   choiceSet,
		UserId:    user_id,
	}

	s.DefaultTemplate(w, r, "question_details.html", temp)
}

func (s *Server) saveVote(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	questionId, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Fatalln("Parsing error")
	}

	// fmt.Println(r.PostForm.Get("choice"))

	var form storage.Choice
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatalln("Decoding error")
	}
	form.Votes += 1
	e := s.store.CreateVote(form)
	//e is refer of error
	if e != nil {
		log.Fatalln("unable to save data: ", err)
	}
	NewUrl := fmt.Sprintf("/get/question/%d/", questionId)
	http.Redirect(w, r, NewUrl, http.StatusSeeOther)
	// fmt.Println("Id: ", id)

	// fmt.Println("=======================Poll===========")
	// fmt.Println(r.ParseForm())

	// fmt.Println(form.ChoiceText)
}
