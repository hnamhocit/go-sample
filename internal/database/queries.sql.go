// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createPost = `-- name: CreatePost :execresult
INSERT INTO
	posts (user_id, title, content)
VALUES
	(?, ?, ?)
`

type CreatePostParams struct {
	UserID  int32  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createPost, arg.UserID, arg.Title, arg.Content)
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO
	users (display_name, email, password)
VALUES
	(?, ?, ?)
`

type CreateUserParams struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser, arg.DisplayName, arg.Email, arg.Password)
}

const deleteMedia = `-- name: DeleteMedia :execresult
DELETE FROM media
where
	id = ?
`

func (q *Queries) DeleteMedia(ctx context.Context, id int32) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteMedia, id)
}

const deletePost = `-- name: DeletePost :execresult
DELETE FROM posts
WHERE
	id = ?
`

func (q *Queries) DeletePost(ctx context.Context, id int32) (sql.Result, error) {
	return q.db.ExecContext(ctx, deletePost, id)
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE
	id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getPost = `-- name: GetPost :one
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
	p.id = ?
`

type GetPostRow struct {
	PostID          int32          `json:"post_id"`
	PostTitle       string         `json:"post_title"`
	PostContent     string         `json:"post_content"`
	PostCreatedAt   time.Time      `json:"post_created_at"`
	PostUpdatedAt   time.Time      `json:"post_updated_at"`
	UserEmail       string         `json:"user_email"`
	UserDisplayName string         `json:"user_display_name"`
	UserPhotoUrl    sql.NullString `json:"user_photo_url"`
}

func (q *Queries) GetPost(ctx context.Context, id int32) (GetPostRow, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i GetPostRow
	err := row.Scan(
		&i.PostID,
		&i.PostTitle,
		&i.PostContent,
		&i.PostCreatedAt,
		&i.PostUpdatedAt,
		&i.UserEmail,
		&i.UserDisplayName,
		&i.UserPhotoUrl,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
SELECT
	id, title, content, created_at, updated_at, user_id
FROM
	posts
WHERE
	id > ?
ORDER BY
	id ASC
LIMIT
	?
`

type GetPostsParams struct {
	ID    int32 `json:"id"`
	Limit int32 `json:"limit"`
}

func (q *Queries) GetPosts(ctx context.Context, arg GetPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT
	id, display_name, email, password, refresh_token, bio, photo_url, background_url, role, token_version, created_at, updated_at
FROM
	users
WHERE
	id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Email,
		&i.Password,
		&i.RefreshToken,
		&i.Bio,
		&i.PhotoUrl,
		&i.BackgroundUrl,
		&i.Role,
		&i.TokenVersion,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
	id, display_name, email, password, refresh_token, bio, photo_url, background_url, role, token_version, created_at, updated_at
FROM
	users
WHERE
	email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.Email,
		&i.Password,
		&i.RefreshToken,
		&i.Bio,
		&i.PhotoUrl,
		&i.BackgroundUrl,
		&i.Role,
		&i.TokenVersion,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserTokenVersion = `-- name: GetUserTokenVersion :one
SELECT
	token_version
FROM
	users
WHERE
	id = ?
`

func (q *Queries) GetUserTokenVersion(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRowContext(ctx, getUserTokenVersion, id)
	var token_version int32
	err := row.Scan(&token_version)
	return token_version, err
}

const getUsers = `-- name: GetUsers :many
SELECT
	id, display_name, email, password, refresh_token, bio, photo_url, background_url, role, token_version, created_at, updated_at
FROM
	users
ORDER BY
	id
LIMIT
	?
OFFSET
	?
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.DisplayName,
			&i.Email,
			&i.Password,
			&i.RefreshToken,
			&i.Bio,
			&i.PhotoUrl,
			&i.BackgroundUrl,
			&i.Role,
			&i.TokenVersion,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePost = `-- name: UpdatePost :execresult
UPDATE posts
SET
	title = ?,
	content = ?
WHERE
	id = ?
`

type UpdatePostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ID      int32  `json:"id"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updatePost, arg.Title, arg.Content, arg.ID)
}

const updateUser = `-- name: UpdateUser :execresult
UPDATE users
SET
	display_name = ?,
	email = ?,
	password = ?,
	bio = ?
WHERE
	id = ?
`

type UpdateUserParams struct {
	DisplayName string         `json:"display_name"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	Bio         sql.NullString `json:"bio"`
	ID          int32          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.DisplayName,
		arg.Email,
		arg.Password,
		arg.Bio,
		arg.ID,
	)
}

const updateUserBackgroundURL = `-- name: UpdateUserBackgroundURL :execresult
UPDATE users
SET
	background_url = ?
WHERE
	id = ?
`

type UpdateUserBackgroundURLParams struct {
	BackgroundUrl sql.NullString `json:"background_url"`
	ID            int32          `json:"id"`
}

func (q *Queries) UpdateUserBackgroundURL(ctx context.Context, arg UpdateUserBackgroundURLParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserBackgroundURL, arg.BackgroundUrl, arg.ID)
}

const updateUserPhotoURL = `-- name: UpdateUserPhotoURL :execresult
UPDATE users
SET
	photo_url = ?
WHERE
	id = ?
`

type UpdateUserPhotoURLParams struct {
	PhotoUrl sql.NullString `json:"photo_url"`
	ID       int32          `json:"id"`
}

func (q *Queries) UpdateUserPhotoURL(ctx context.Context, arg UpdateUserPhotoURLParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserPhotoURL, arg.PhotoUrl, arg.ID)
}

const updateUserRefreshToken = `-- name: UpdateUserRefreshToken :execresult
UPDATE users
SET
	refresh_token = ?,
	token_version = token_version + 1
WHERE
	id = ?
`

type UpdateUserRefreshTokenParams struct {
	RefreshToken sql.NullString `json:"refresh_token"`
	ID           int32          `json:"id"`
}

func (q *Queries) UpdateUserRefreshToken(ctx context.Context, arg UpdateUserRefreshTokenParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserRefreshToken, arg.RefreshToken, arg.ID)
}

const uploadMedia = `-- name: UploadMedia :execresult
INSERT INTO
	media (user_id, name, content_type, path, size)
VALUES
	(?, ?, ?, ?, ?)
`

type UploadMediaParams struct {
	UserID      int32  `json:"user_id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Path        string `json:"path"`
	Size        int32  `json:"size"`
}

func (q *Queries) UploadMedia(ctx context.Context, arg UploadMediaParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, uploadMedia,
		arg.UserID,
		arg.Name,
		arg.ContentType,
		arg.Path,
		arg.Size,
	)
}
