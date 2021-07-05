package handler

import (
	"net/http"
)

func (s *Server) logOut(w http.ResponseWriter, r *http.Request) {
	session, _ := s.session.Get(r, "poll_app")
	session.Values["user_email"] = nil
	session.Values["user_id"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
