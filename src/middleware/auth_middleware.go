package middleware

import (
	"log"
	"os"
	"time"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var identityKey = "id"

// AuthMiddleware returns middleware for authentication using JWT authorization
func AuthMiddleware() *jwt.GinJWTMiddleware {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "reservation zone",
		Key:           []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   identityKey,
		PayloadFunc:   convertStructToMapClaims,
		Authenticator: verifyCredential,
		Authorizator:  isAuthorized,
		Unauthorized:  failedAuthorization,
		TokenLookup:   "header: Authorization, query: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

func convertStructToMapClaims(data interface{}) jwt.MapClaims {
	if v, ok := data.(*entity.UserAuth); ok {
		return jwt.MapClaims{
			identityKey: v.Email,
		}
	}
	return jwt.MapClaims{}
}

func handleIdentity(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &entity.User{
		Email: claims[identityKey].(string),
	}
}

func verifyCredential(c *gin.Context) (interface{}, error) {
	var input entity.UserInput
	if err := c.ShouldBind(&input); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	if input.UserName == "" {
		return login(input)
	}
	return signin(input)
}

func signin(input entity.UserInput) (*entity.UserAuth, error) {
	username := input.UserName
	email := input.Email

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return &entity.UserAuth{}, err
	}

	var us service.UserService
	u, err := us.CreateModel(username, email, hashedPassword)
	if err != nil {
		return &entity.UserAuth{}, err
	}

	return &entity.UserAuth{
		ID:       u.ID,
		Email:    email,
		Password: hashedPassword,
	}, nil
}

func login(input entity.UserInput) (p *entity.UserAuth, err error) {
	email := input.Email
	password := input.Password

	var us service.UserService
	storedUser, err := us.FindByEmail(email)
	if err != nil {
		return &entity.UserAuth{}, err
	}

	if err := bcrypt.CompareHashAndPassword(storedUser.Password, []byte(password)); err != nil {
		return &entity.UserAuth{}, jwt.ErrFailedAuthentication
	}

	return &entity.UserAuth{
		ID:       storedUser.ID,
		Email:    storedUser.Email,
		Password: storedUser.Password,
	}, nil
}

func isAuthorized(data interface{}, c *gin.Context) bool {
	var ID uint
	if err := c.ShouldBindUri(&ID); err != nil {
		return false
	}
	if v, ok := data.(*entity.UserAuth); ok && v.ID == ID {
		return true
	}
	return false
}

func failedAuthorization(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
