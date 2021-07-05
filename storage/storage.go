package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	User struct {
		ID        int32     `db:"id"`
		Name      string    `db:"name"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Password  string    `db:"password"`
		IsActive  bool      `db:"is_active"`
		IsAdmin   bool      `db:"is_admin"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	Question struct {
		ID           int32     `db:"id"`
		UserID       int32     `db:"user_id"`
		QuestionText string    `db:"question_text"`
		CreatedAt    time.Time `db:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"`
	}
	Choice struct {
		ID         int32     `db:"id"`
		QuestionID int32     `db:"question_id"`
		UserID     int32     `db:"user_id"`
		ChoiceText string    `db:"choice_text"`
		Votes      int32     `db:"votes"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}
)

func (sg User) Validate() error {
	return validation.ValidateStruct(&sg,
		validation.Field(&sg.Name,
			validation.Required.Error("name is required"),
			validation.Length(5, 100).Error("name length must be 5 to 100"),
		),
		validation.Field(&sg.Username,
			validation.Required.Error("username is required"),
			validation.Length(3, 20).Error("usrname length must be 3 to 20"),
		),
		validation.Field(&sg.Email,
			validation.Required.Error("email is required"),
			is.Email,
		),
		validation.Field(&sg.Password,
			validation.Required.Error("Password is required"),
			validation.Length(3, 10).Error("Password Lenght must be 3 to 10"),
		),
	)
}

func (sg User) ValidateUser() error {
	return validation.ValidateStruct(&sg,
		validation.Field(&sg.Email,
			validation.Required.Error("email is required"),
			is.Email,
		),
		validation.Field(&sg.Password,
			validation.Required.Error("Password is required"),
			validation.Length(3, 10).Error("Password Lenght must be 3 to 10"),
		),
	)
}

func (q Question) ValidateQuestion() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.QuestionText,
			validation.Required.Error("question is required"),
			validation.Length(10, 350).Error("Length must be 10 to 350"),
		),
	)
}

func (c Choice) ValidateChoice() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ChoiceText,
			validation.Required.Error("Choice is required"),
			validation.Length(2, 350).Error("Choice text length must be 2 to 350"),
		),
	)
}
