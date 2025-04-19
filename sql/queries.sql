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

-- name: UploadMedia :execresult
INSERT INTO
	media (user_id, name, content_type, path, size)
VALUES
	(?, ?, ?, ?, ?);

-- name: DeleteMedia :execresult
DELETE FROM media
where
	id = ?;

-- name: GetPosts :many
SELECT
	*
FROM
	posts
WHERE
	id > ?
ORDER BY
	id ASC
LIMIT
	?;

-- name: GetPost :one
SELECT
	p.id AS post_id,
	p.title AS post_title,
	p.content AS post_content,
	p.created_at AS post_created_at,
	p.updated_at AS post_updated_at,
	u.email AS user_email,
	u.display_name AS user_display_name,
	u.photo_url AS user_photo_url
FROM
	posts p
	JOIN users u ON p.user_id = u.id
WHERE
	p.id = ?;

-- name: DeletePost :execresult
DELETE FROM posts
WHERE
	id = ?;

-- name: CreatePost :execresult
INSERT INTO
	posts (user_id, title, content)
VALUES
	(?, ?, ?);

-- name: UpdatePost :execresult
UPDATE posts
SET
	title = ?,
	content = ?
WHERE
	id = ?;
