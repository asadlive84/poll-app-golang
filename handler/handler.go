package handler

import (
	"html/template"
	"log"
	"net/http"
	"poll-app/storage/postgres"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

type (
	Server struct {
		templates *template.Template
		store     *postgres.Storage
		decoder   *schema.Decoder
		session   *sessions.CookieStore
	}
)

func NewServer(st *postgres.Storage, decoder *schema.Decoder, session *sessions.CookieStore) (*mux.Router, error) {

	s := &Server{
		templates: &template.Template{},
		store:     st,
		decoder:   decoder,
		session:   session,
	}

	if err := s.parseTemplates(); err != nil {
		log.Println("parse template error")
	}

	r := mux.NewRouter()
	r.Use(csrf.Protect([]byte("1234")))

	r.HandleFunc("/", s.homePage)

	r.HandleFunc("/login/", s.getLoginPage).Methods("GET")
	r.HandleFunc("/login/", s.postLoginPage).Methods("POST")
	r.HandleFunc("/logout/", s.logOut)
	r.HandleFunc("/signup/", s.getSignupPage).Methods("GET")
	r.HandleFunc("/signup/", s.postSignupPage).Methods("POST")

	ar := r.NewRoute().Subrouter()

	ar.Use(s.AuthMiddleware)

	ar.HandleFunc("/create/question/", s.getCreateQuesionPage).Methods("GET")
	ar.HandleFunc("/create/question/", s.postCreateQuesionPage).Methods("POST")

	ar.HandleFunc("/question/create/choice/{id}/", s.getCreateChoicePage).Methods("GET")
	ar.HandleFunc("/question/create/choice/{id}/", s.postCreateChoicePage).Methods("POST")

	ar.HandleFunc("/question/update/choice/{id}/{cId}/", s.getUpdateChoicePage).Methods("GET")
	ar.HandleFunc("/question/update/choice/{id}/{cId}/", s.postUpdateChoicePage).Methods("POST")

	ar.HandleFunc("/get/question/", s.getAllQuestionPage).Methods("GET")

	ar.HandleFunc("/get/question/{id}/", s.questionDetails).Methods("GET")

	ar.HandleFunc("/update/question/{id}/", s.getUpdateQuestionPage).Methods("GET")
	ar.HandleFunc("/update/question/{id}/", s.postUpdateQuestionPage).Methods("POST")

	ar.HandleFunc("/get/question/vote/{id}/", s.saveVote).Methods("POST")

	return r, nil

}

func (s *Server) parseTemplates() error {
	templates := template.New("templates").Funcs(template.FuncMap{
		"strrev": func(str string) string {
			n := len(str)
			runes := make([]rune, n)
			for _, rune := range str {
				n--
				runes[n] = rune
			}
			return string(runes[n:])
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseGlob("assets/templates/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		session, _ := s.session.Get(r, "poll_app")
		value := session.Values["user_id"]
		user_email := session.Values["user_email"]
		log.Println(value)
		if _, ok := user_email.(string); ok {
			log.Println(r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden, You're not regirsted", http.StatusForbidden)
		}

	})
}

func (s *Server) DefaultTemplate(w http.ResponseWriter, r *http.Request, temp_name string, data interface{}) {
	temp := s.templates.Lookup(temp_name)

	if err := temp.Execute(w, data); err != nil {
		log.Fatalln("executing template: ", err)
		return
	}
}
