package postgres

import (
	"log"
	"poll-app/storage"
)

const createQuestionQuery = `
	INSERT INTO question 
	(user_id, question_text) 
	VALUES (:user_id, :question_text) 
	RETURNING id, created_at, updated_at
`

func (s *Storage) SaveQuestionDB(question storage.Question) (int32, error) {
	stmt, err := s.db.PrepareNamed(createQuestionQuery)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, question); err != nil {
		return 0, err
	}

	return id, nil

}

const updateQuesionQuery = `
		UPDATE question 
		SET question_text=$1,user_id=$2 WHERE id=$3 RETURNING *
`

func (s *Storage) UpdateQuestionDB(question_text string, user_id int32, id int32) (storage.Question, error) {
	question := storage.Question{}

	err := s.db.Get(&question, updateQuesionQuery, question_text, user_id, id)
	return question, err
}

const getAllQuestion = `SELECT * FROM question ORDER BY id DESC`

func (s *Storage) GetAllQuestionDB() ([]storage.Question, error) {
	question := make([]storage.Question, 0)
	if err := s.db.Select(&question, getAllQuestion); err != nil {
		return nil, err
	}

	return question, nil
}

const getSingleQuestionQuery = `SELECT * FROM question WHERE id=$1 limit 1`

func (s *Storage) GetQuestionDetail(id int32) (*storage.Question, error) {
	question := storage.Question{}
	if err := s.db.Get(&question, getSingleQuestionQuery, id); err != nil {
		return nil, err
	}
	return &question, nil

}
