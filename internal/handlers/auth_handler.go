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

type AuthRepo struct {
	Dao *database.Queries
	Ctx context.Context
}

func (r *AuthRepo) UpdateRefreshToken(id int32, refreshToken *string) error {
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

	_, err := r.Dao.UpdateUser(r.Ctx, database.UpdateUserParams{
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

func (r *AuthRepo) Login(c *gin.Context) {
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

	tokens, err := utils.GenerateTokens(existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	updateRefreshTokenErr := r.UpdateRefreshToken(existingUser.ID, &tokens.RefreshToken)
	if updateRefreshTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": updateRefreshTokenErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Login successfuly!",
		"data": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}

type RegisterDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	DisplayName string `json:"displayName" binding:"required,min=1,max=35"`
}

func (r *AuthRepo) Register(c *gin.Context) {
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

	tokens, err := utils.GenerateTokens(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	updateRefreshTokenErr := r.UpdateRefreshToken(int32(id), &tokens.RefreshToken)
	if updateRefreshTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  updateRefreshTokenErr.Error(),
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

func (r *AuthRepo) Logout(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Invalid user ID",
		})
		return
	}

	err := r.UpdateRefreshToken(int32(userId.(int)), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Logout successfully!",
	})
}

func (r *AuthRepo) Refresh(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Invalid user ID",
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

	user, err := r.Dao.GetUser(r.Ctx, int32(userId.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	isValidRefreshToken, err := utils.Verify(refreshToken.(string), user.RefreshToken.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	if !isValidRefreshToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "Invalid refresh token",
		})
		return
	}

	tokens, err := utils.GenerateTokens(int32(userId.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	updateRefreshTokenErr := r.UpdateRefreshToken(int32(userId.(int)), &tokens.RefreshToken)
	if updateRefreshTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  updateRefreshTokenErr.Error(),
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
