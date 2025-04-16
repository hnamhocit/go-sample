-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN token_version INT DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN token_version;

-- +goose StatementEnd
