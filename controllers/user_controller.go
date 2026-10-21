package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(context *gin.Context) {
	var input AuthInputRegister

	//Validation
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Your input is not valid",
			"error":   err.Error(),
		})
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to encrypt password",
			"error":   errHash.Error(),
		})
		return
	}

	// save to DB
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to registration",
			"error":   userCreated.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully user registration",
		"data": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})
}

func LoginUser(context *gin.Context) {
	var input AuthInputLogin

	//Validation
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Your input is not valid",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	userData := config.DB.Where("email = ?", input.Email).First(&user).Error
	if userData != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email is not Registered",
			"error":   userData.Error(),
		})
		return
	}

	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatchPassword != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Your password is not valid",
			"error":   errMatchPassword.Error(),
		})
		return
	}

	// Buat Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	// Jwt secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to created token",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully Login",
		"token":   tokenString,
		"data": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})
}

func GetCurrentUser(context *gin.Context) {
	userId, exists := context.Get("userID")

	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not Authentication",
		})
		return
	}

	var user models.User
	userData := config.DB.Select("id", "name", "email").First(&user, userId).Error
	if userData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "User is not found!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully to get user data",
		"data": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})
	return
}
