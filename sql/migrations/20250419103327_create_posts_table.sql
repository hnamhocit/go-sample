-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
	id INT PRIMARY KEY AUTO_INCREMENT,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created_at DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	user_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;

-- +goose StatementEnd
