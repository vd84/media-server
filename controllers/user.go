package controllers

import (
	"mediaserver/data"
	"mediaserver/views"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterAccountEndpoints(e *gin.Engine) {
	e.POST("/account/createUser", createUser)
	e.POST("/account/login", login)
	e.POST("/account/logout", logout)
}

func createUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var userView views.User
	if err := c.ShouldBindJSON(&userView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}

	var userToSave data.User
	if err := db.Where("username = ?", userView.Username).First(&userToSave).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username or email already exists",
		})
		return
	}

	hashedPassword, err := HashPassword(userView.Password)
	if err != nil {
		panic("failed to hash password")
	}

	userToSave.HashedPassword = hashedPassword
	userToSave.Username = userView.Username

	if err := db.Create(&userToSave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userToSave.ID,
		"message": "User created successfully",
	})
}

func login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var userView views.User
	if err := c.ShouldBindJSON(&userView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var userDb data.User
	if err := db.Where(&data.User{Username: userView.Username}).First(&userDb).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.HashedPassword), []byte(userView.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := createToken(userDb.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": userDb.ID})
}

func logout(c *gin.Context) {

	c.SetCookie("token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func createToken(userID uint) (string, error) {
	jwtSecret := []byte(os.Getenv("JWTSECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString(jwtSecret)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
