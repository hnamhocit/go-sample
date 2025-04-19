package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sample/internal/database"
	"sample/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Dao *database.Queries
	Ctx context.Context
}

func (r *AuthHandler) UpdateRefreshToken(id int32, refreshToken *string) error {
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

	_, err := r.Dao.UpdateUserRefreshToken(r.Ctx, database.UpdateUserRefreshTokenParams{
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

func (r *AuthHandler) Login(c *gin.Context) {
	var input LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	existingUser, err := r.Dao.GetUserByEmail(r.Ctx, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	ok, err := utils.Verify(input.Password, existingUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": "Invalid credentials"})
		return
	}

	tokens, err := utils.GenerateTokens(existingUser.ID, existingUser.TokenVersion+1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	updateErr := r.UpdateRefreshToken(existingUser.ID, nil)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  updateErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Login successfully!",
		"data": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}

type RegisterDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=35"`
}

func (r *AuthHandler) Register(c *gin.Context) {
	var input RegisterDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	_, err := r.Dao.GetUserByEmail(r.Ctx, input.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{
			"code": 0,
			"msg":  "Email already exists!",
		})
		return
	}

	hashedPassword, err := utils.Hash(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	newUser, err := r.Dao.CreateUser(r.Ctx, database.CreateUserParams{
		DisplayName: input.DisplayName,
		Email:       input.Email,
		Password:    hashedPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	id, err := newUser.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	tokens, err := utils.GenerateTokens(int32(id), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	updateErr := r.UpdateRefreshToken(int32(id), nil)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  updateErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Register successfully!",
		"data": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}

func (r *AuthHandler) Logout(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Invalid user ID",
		})
		return
	}

	id, ok := userId.(int32)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Covert user ID error",
		})
		return
	}

	updateErr := r.UpdateRefreshToken(id, nil)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  updateErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Logout successfully!",
	})
}

func (r *AuthHandler) Refresh(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Invalid user ID",
		})
		return
	}

	id, ok := userId.(int32)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Covert user ID error",
		})
		return
	}

	refreshToken, ok := c.Get("refresh_token")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Invalid refresh token",
		})
		return
	}

	user, err := r.Dao.GetUser(r.Ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	ok, verifyErr := utils.Verify(refreshToken.(string), user.RefreshToken.String)
	if verifyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  verifyErr.Error(),
		})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 0,
			"msg":  "Invalid refresh token!",
		})
		return
	}

	tokens, err := utils.GenerateTokens(id, user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	updateErr := r.UpdateRefreshToken(id, &tokens.RefreshToken)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  updateErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Refresh token successfully!",
		"data": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}
