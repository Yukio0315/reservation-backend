package middleware

import (
	"log"
	"os"
	"time"

	"github.com/Yukio0315/reservation-backend/src/api"
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
		Realm:           "reservation zone",
		Key:             []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     convertUserIDToMapClaims,
		IdentityHandler: handleIdentity,
		Authenticator:   verifyCredential,
		Authorizator:    isAuthorized,
		Unauthorized:    failedAuthorization,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

func convertUserIDToMapClaims(data interface{}) jwt.MapClaims {
	if v, ok := data.(*entity.UserIDAndPassword); ok {
		return jwt.MapClaims{
			identityKey: v.ID,
		}
	}
	return jwt.MapClaims{}
}

func handleIdentity(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return claims[identityKey].(float64)
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

func login(input entity.UserInput) (p *entity.UserIDAndPassword, err error) {
	var us service.UserService
	storedUser, err := us.FindIDAndPasswordByEmail(input.Email)
	if err != nil {
		return &entity.UserIDAndPassword{}, err
	}

	if err := bcrypt.CompareHashAndPassword(storedUser.Password, []byte(input.Password)); err != nil {
		return &entity.UserIDAndPassword{}, jwt.ErrFailedAuthentication
	}

	return &entity.UserIDAndPassword{
		ID:       storedUser.ID,
		Password: storedUser.Password,
	}, nil
}

func signin(input entity.UserInput) (*entity.UserIDAndPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return &entity.UserIDAndPassword{}, err
	}

	var us service.UserService
	u, err := us.CreateModel(input.UserName, input.Email, hashedPassword)
	if err != nil {
		return &entity.UserIDAndPassword{}, err
	}

	go api.GmailContent{
		Email:   input.Email,
		Subject: "【シェアオフィス】ご登録ありがとうございます",
		Body:    "シェアオフィスへのご登録が完了しました。",
	}.Send()

	return &entity.UserIDAndPassword{
		ID:       u.ID,
		Password: hashedPassword,
	}, nil
}

func isAuthorized(data interface{}, c *gin.Context) bool {
	var id entity.UserID
	if err := c.ShouldBindUri(&id); err != nil {
		return false
	}
	if v, ok := data.(float64); ok && v == float64(id.ID) {
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
