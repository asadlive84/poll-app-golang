package handler

import (
	"log"
	"net/http"
)

type (
	HomePage struct {
		UserLoggedIn bool
	}
)

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	temp := s.templates.Lookup("home.html")
	if temp == nil {
		log.Fatalln("template not loading")
	}

	session, _ := s.session.Get(r, "poll_app")
	userID := session.Values["user_id"]

	if _, ok := userID.(string); ok {
		data := HomePage{
			UserLoggedIn: true,
		}
		if err := temp.Execute(w, data); err != nil {
			log.Fatalln("Session Execution error")
		}
		return
	}

	data := HomePage{
		UserLoggedIn: false,
	}

	if err := temp.Execute(w, data); err != nil {
		log.Fatalln("temp Execution error")
	}

}
