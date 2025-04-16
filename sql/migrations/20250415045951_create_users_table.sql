-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	id INT PRIMARY KEY AUTO_INCREMENT,
	display_name VARCHAR(35) NOT NULL,
	email VARCHAR(100) NOT NULL UNIQUE,
	password TEXT NOT NULL,
	refresh_token TEXT DEFAULT NULL,
	bio TEXT DEFAULT NULL,
	photo_url TEXT DEFAULT NULL,
	background_url TEXT DEFAULT NULL,
	role ENUM ('USER', 'ADMIN') DEFAULT 'USER',
	created_at DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
