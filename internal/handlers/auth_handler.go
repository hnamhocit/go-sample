package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sample/internal/database"
	"sample/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler
}

func (h *AuthHandler) UpdateRefreshToken(id int32, refreshToken *string) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID: %d", id)
	}

	var hashedRefreshToken sql.NullString
	if refreshToken != nil {
		hashed, err := utils.Hash(*refreshToken)
		if err != nil {
			return err
		}

		hashedRefreshToken = sql.NullString{
			String: hashed,
			Valid:  true,
		}
	} else {
		hashedRefreshToken = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	_, err := h.Dao.UpdateUserRefreshToken(h.Ctx, database.UpdateUserRefreshTokenParams{
		ID:           id,
		RefreshToken: hashedRefreshToken,
	})
	if err != nil {
		return err
	}

	return nil
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, err.Error())
		return
	}

	existingUser, err := h.Dao.GetUserByEmail(h.Ctx, input.Email)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	ok, err := utils.Verify(input.Password, existingUser.Password)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	if !ok {
		h.handleError(c, "Invalid credentials")
		return
	}

	tokens, err := utils.GenerateTokens(existingUser.ID, existingUser.TokenVersion+1)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	updateErr := h.UpdateRefreshToken(existingUser.ID, nil)
	if updateErr != nil {
		h.handleError(c, updateErr.Error())
		return
	}

	h.handleSuccess(c, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}, nil)
}

type RegisterDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=35"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	_, err := h.Dao.GetUserByEmail(h.Ctx, input.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.handleError(c, err.Error())
			return
		}
	} else {
		h.handleError(c, "Email already exists!")
		return
	}

	hashedPassword, err := utils.Hash(input.Password)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	newUser, err := h.Dao.CreateUser(h.Ctx, database.CreateUserParams{
		DisplayName: input.DisplayName,
		Email:       input.Email,
		Password:    hashedPassword,
	})
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	id, err := newUser.LastInsertId()
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	tokens, err := utils.GenerateTokens(int32(id), 0)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	updateErr := h.UpdateRefreshToken(int32(id), nil)
	if updateErr != nil {
		h.handleError(c, updateErr.Error())
		return
	}

	h.handleSuccess(c, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	id, ok := utils.GetUserID(c)
	if !ok {
		h.handleError(c, "Invalid user ID")
		return
	}

	updateErr := h.UpdateRefreshToken(*id, nil)
	if updateErr != nil {
		h.handleError(c, updateErr.Error())
		return
	}

	h.handleSuccess(c, nil, nil)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	id, ok := utils.GetUserID(c)
	if !ok {
		h.handleError(c, "Invalid user ID")
		return
	}

	refreshToken, ok := c.Get("refresh_token")
	if !ok {
		h.handleError(c, "Invalid refresh token")
		return
	}

	user, err := h.Dao.GetUser(h.Ctx, *id)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	ok, verifyErr := utils.Verify(refreshToken.(string), user.RefreshToken.String)
	if verifyErr != nil {
		h.handleError(c, verifyErr.Error())
		return
	}

	if !ok {
		h.handleError(c, "Invalid refresh token!")
		return
	}

	tokens, err := utils.GenerateTokens(*id, user.TokenVersion)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	updateErr := h.UpdateRefreshToken(*id, &tokens.RefreshToken)
	if updateErr != nil {
		h.handleError(c, updateErr.Error())
		return
	}

	h.handleSuccess(c, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}, nil)
}
