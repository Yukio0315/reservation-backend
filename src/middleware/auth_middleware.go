package middleware

import (
	"log"
	"os"
	"time"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/Yukio0315/reservation-backend/src/util"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

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
		IdentityKey:     util.IDENTITYKEY,
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
	if v, ok := data.(*entity.UserAuth); ok {
		return jwt.MapClaims{
			util.IDENTITYKEY:  v.ID,
			util.IDENTITYKEY2: v.Permission,
		}
	}
	return jwt.MapClaims{}
}

func handleIdentity(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return claims[util.IDENTITYKEY].(float64)
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

func login(input entity.UserInput) (*entity.UserAuth, error) {
	us := service.UserService{}
	storedUser, err := us.FindByEmail(input.Email)
	if err != nil {
		return &entity.UserAuth{}, err
	}

	if err = bcrypt.CompareHashAndPassword(storedUser.Password, []byte(input.Password)); err != nil {
		return &entity.UserAuth{}, jwt.ErrFailedAuthentication
	}

	return &entity.UserAuth{
		ID:         storedUser.ID,
		Password:   storedUser.Password,
		Permission: storedUser.Permission,
	}, nil
}

func signin(input entity.UserInput) (*entity.UserAuth, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return &entity.UserAuth{}, err
	}

	us := service.UserService{}
	u, err := us.CreateModel(input.UserName, input.Email, hashedPassword)
	if err != nil {
		return &entity.UserAuth{}, err
	}

	go api.GmailContent{
		Email:   input.Email,
		Subject: "【シェアオフィス】ご登録ありがとうございます",
		Body:    "シェアオフィスへのご登録が完了しました。",
	}.Send()

	return &entity.UserAuth{
		ID:         u.ID,
		Permission: u.Permission,
		Password:   hashedPassword,
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

// AdminMiddleware check administrator authorization
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		if claims[util.IDENTITYKEY2].(string) != "admin" {
			failedAuthorization(c, 400, "Invalid token")
			return
		}
		c.Next()
	}
}
