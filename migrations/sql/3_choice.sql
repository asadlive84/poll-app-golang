-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS choice(
   id INT GENERATED ALWAYS AS IDENTITY,
   question_id INT DEFAULT NULL,
   user_id INT DEFAULT NULL,
   choice_text VARCHAR(255) NOT NULL,
   votes INT default 0,
   created_at timestamp default current_timestamp,
   updated_at timestamp default current_timestamp,
   PRIMARY KEY(id),
   CONSTRAINT question_id FOREIGN KEY(question_id) REFERENCES question(id),
   CONSTRAINT user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS choice;
-- +goose StatementEnd