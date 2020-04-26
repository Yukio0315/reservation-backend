package controller

import (
	"errors"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/Yukio0315/reservation-backend/src/template"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserController type
type UserController struct {
	us service.UserService
	os service.OneTimeURLService
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

	u, err := uc.us.FindByID(id.ID)
	if err != nil {
		c.AbortWithError(404, err)
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
		Subject: template.CHANGEPASSWORDTITLE,
		Body:    template.CHANGEPASSWORDBODY,
	}.Send()
}

// ReserveResetPassword is reservation for reset password. It create one time url and send email.
func (uc UserController) ReserveResetPassword(c *gin.Context) {
	input := entity.UserEmail{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(400, err)
		return
	}

	u, err := uc.us.FindByEmail(input.Email)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	o, err := uc.os.Create(u.ID)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Status(200)

	go api.GmailContent{
		Email:   input.Email,
		Subject: template.ONETIMEURLTITLE,
		Body:    template.OneTimeURLBody(o.QueryString),
	}.Send()
}

// PasswordReset controls resetting password
func (uc UserController) PasswordReset(c *gin.Context) {
	input := entity.UserInputMailPassword{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(400, err)
		return
	}
	query := entity.OneTimeQuery{}
	if err := c.ShouldBindUri(&query); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := uc.os.DeleteByUUID(query.UUID); err != nil {
		c.AbortWithError(400, err)
		return
	}

	u, err := uc.us.FindByEmail(input.Email)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	if err := uc.us.UpdatePassword(u.ID, input.Password); err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)

	go api.GmailContent{
		Email:   input.Email,
		Subject: template.RESETPASSWORDTITLE,
		Body:    template.RESETPASSWORDBODY,
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
	c.JSON(200, entity.UserEmail{
		Email: input.Email,
	})

	go api.GmailContent{
		Email:   input.Email,
		Subject: template.CHANGEEMAILTITLE,
		Body:    template.CHANGEEMAILBODY,
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

	u, err := uc.us.FindByEmail(input.Email)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	if id.ID != u.ID {
		c.AbortWithError(400, errors.New("invalid email"))
		return
	}

	if err := uc.us.DeleteByID(id.ID); err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Status(204)

	go api.GmailContent{
		Email:   input.Email,
		Subject: template.DELETEACCOUNTTITLE,
		Body:    template.DELETEACCOUNTBODY,
	}.Send()
}
