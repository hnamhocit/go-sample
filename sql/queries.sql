-- name: GetUsers :many
SELECT
	*
FROM
	users
ORDER BY
	id
LIMIT
	?
OFFSET
	?;

-- name: GetUserTokenVersion :one
SELECT
	token_version
FROM
	users
WHERE
	id = ?;

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
	display_name = ?,
	email = ?,
	password = ?,
	bio = ?
WHERE
	id = ?;

-- name: UpdateUserRefreshToken :execresult
UPDATE users
SET
	refresh_token = ?,
	token_version = token_version + 1
WHERE
	id = ?;

-- name: UpdateUserPhotoURL :execresult
UPDATE users
SET
	photo_url = ?
WHERE
	id = ?;

-- name: UpdateUserBackgroundURL :execresult
UPDATE users
SET
	background_url = ?
WHERE
	id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
	id = ?;
