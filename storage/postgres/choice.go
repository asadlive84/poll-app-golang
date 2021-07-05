package postgres

import (
	"log"
	"poll-app/storage"
)

const createChoiceQuery = `
	INSERT INTO choice
	(question_id, user_id,choice_text) VALUES (:question_id, :user_id, :choice_text)
	RETURNING id, created_at, updated_at
`

func (s *Storage) SaveChoiceDB(choice storage.Choice) (int32, error) {
	stmt, err := s.db.PrepareNamed(createChoiceQuery)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, choice); err != nil {
		return 0, err
	}

	return id, nil

}

const ChoiceQuery = `select * from choice where question_id=$1 ORDER BY id desc`

func (s *Storage) ChoiceQuery(id int32) ([]storage.Choice, error) {
	choice := make([]storage.Choice, 0)
	if err := s.db.Select(&choice, ChoiceQuery, id); err != nil {
		return nil, err
	}
	return choice, nil
}

const createPollVote = `UPDATE choice SET votes=$1 WHERE id=$2;`

func (s *Storage) CreateVote(choice storage.Choice) error {
	_, err := s.db.Exec(createPollVote, choice.Votes, choice.ID)

	if err != nil {
		return err
	}

	return nil
}

const updateChoiceQuery = `
		UPDATE choice 
		SET choice_text=$1 WHERE id=$2 RETURNING *
`

func (s *Storage) UpdateChoiceDB(choice_text string, id int32) (storage.Choice, error) {
	log.Println("===============")
	log.Println(choice_text)
	log.Println(id)
	log.Println("================")
	choice := storage.Choice{}

	err := s.db.Get(&choice, updateChoiceQuery, choice_text, id)
	return choice, err
}

const getSingleChoiceQuery = `SELECT * FROM choice WHERE id=$1 AND question_id=$2  limit 1`

func (s *Storage) GetChoiceDetail(id, question_id int32) (*storage.Choice, error) {
	choice := storage.Choice{}
	if err := s.db.Get(&choice, getSingleChoiceQuery, id, question_id); err != nil {
		return nil, err
	}
	return &choice, nil

}
