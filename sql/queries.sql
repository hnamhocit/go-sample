-- name: GetUsers :many
SELECT
	*
FROM
	users;

-- name: GetUser :one
SELECT
	*
FROM
	users
WHERE
	id = ?;

-- name: GetUserByEmail :one
SELECT
	*
FROM
	users
WHERE
	email = ?;

-- name: CreateUser :execresult
INSERT INTO
	users (display_name, email, password)
VALUES
	(?, ?, ?);

-- name: UpdateUser :execresult
UPDATE users
SET
	display_name = COALESCE(?, display_name),
	email = COALESCE(?, email),
	password = COALESCE(?, password),
	refresh_token = COALESCE(?, refresh_token)
WHERE
	id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
	id = ?;
