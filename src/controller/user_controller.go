package controller

import (
	"errors"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//TODO: change status code

// UserController type
type UserController struct {
	us service.UserService
}

// Show controlles user information & reservation
func (uc UserController) Show(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(400, err)
		return
	}

	p, err := uc.us.FindUserProfileByID(id.ID)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.JSON(200, p)
}

// PasswordChange controls updating password
func (uc UserController) PasswordChange(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(400, err)
		return
	}

	passwords := entity.UserNewOldPasswords{}
	if err := c.ShouldBindJSON(&passwords); err != nil {
		c.AbortWithError(400, err)
		return
	}

	u, err := uc.us.FindEmailAndPasswordByID(id.ID)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(passwords.OldPassword)); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := uc.us.UpdatePassword(id.ID, passwords.NewPassword); err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)

	go api.GmailContent{
		Email:   u.Email,
		Subject: "【シェアオフィス】パスワードの変更が完了しました",
		Body:    "パスワードの変更が完了しました。",
	}.Send()
}

// PasswordReset controls resetting password
func (uc UserController) PasswordReset(c *gin.Context) {
	input := entity.UserInputMailPassword{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(400, err)
		return
	}

	id, err := uc.us.FindIDByEmail(input.Email)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := uc.us.UpdatePassword(id, input.Password); err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)

	go api.GmailContent{
		Email:   input.Email,
		Subject: "【シェアオフィス】パスワードのリセットが完了しました",
		Body:    "パスワードのリセットが完了しました。",
	}.Send()

}

// UserNameChange chaneg the user name
func (uc UserController) UserNameChange(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(400, err)
		return
	}

	input := entity.UserNameInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := uc.us.UpdateUserNameByID(id.ID, input.UserName); err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Status(200)
}

// EmailChange change the user email
func (uc UserController) EmailChange(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(400, err)
		return
	}

	input := entity.UserEmail{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := uc.us.UpdateEmailByID(id.ID, input.Email); err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Status(200)

	go api.GmailContent{
		Email:   input.Email,
		Subject: "【シェアオフィス】Emailアドレスを変更しました",
		Body:    "Emailアドレスを変更しました。",
	}.Send()
}

// Delete delete the user account
func (uc UserController) Delete(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(400, err)
		return
	}

	input := entity.UserEmail{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(400, err)
		return
	}

	storedID, err := uc.us.FindIDByEmail(input.Email)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	if id.ID != storedID {
		c.AbortWithError(400, errors.New("invalid email"))
		return
	}

	if err := uc.us.DeleteByID(id.ID); err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Status(200)

	go api.GmailContent{
		Email:   input.Email,
		Subject: "【シェアオフィス】アカウントを削除しました",
		Body:    "アカウントを削除しました。\nまたのご利用をお待ちしております。",
	}.Send()
}
