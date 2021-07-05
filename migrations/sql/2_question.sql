-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS question(
   id INT GENERATED ALWAYS AS IDENTITY,
   user_id INT,
   question_text VARCHAR(255) NOT NULL,
   created_at timestamp default current_timestamp,
   updated_at timestamp default current_timestamp,
   PRIMARY KEY(id),
   CONSTRAINT fk_users
      FOREIGN KEY(user_id)
	  REFERENCES users(id)
	  ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS question;
-- +goose StatementEnd