-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS media (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(100) NOT NULL,
	path TEXT NOT NULL,
	content_type VARCHAR(100) NOT NULL,
	size INT NOT NULL,
	created_at DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	user_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS media;

-- +goose StatementEnd
