package main

import (
	"fmt"
	"log"
	"mediaserver/controllers"
	db2 "mediaserver/db"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func main() {
	if err := controllers.EnsureMovieDirExists(); err != nil {
		log.Fatal("Could not create movie directory:", err)
	}
	local := os.Getenv("LOCALDB") == "true"
	var db *gorm.DB
	var err error
	if local {
		db, err = db2.ConnectToLocalDb()
	} else {
		db, err = db2.ConnectToDb()
	}
	db2.RunMigration(db)
	if err != nil {
		return
	}
	r := gin.Default()

	r.Use()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	r.Use(CORSMiddleware())
	r.Use(jwtMiddleware())

	controllers.RegisterAccountEndpoints(r)
	controllers.RegisterMovieEndpoints(r)

	fmt.Println("Starting media server on http://localhost:8080...")
	err = r.Run(":8080")
	if err != nil {
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept,X-Requested-With, X-Filename")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
func jwtMiddleware() gin.HandlerFunc {
	jwtSecret := []byte(os.Getenv("JWTSECRET"))

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/account") {
			c.Next()
			return
		}

		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)

		c.Next()
	}
}
