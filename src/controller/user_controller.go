package controller

import (
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//TODO: change status code

// UserController type
type UserController struct {
	s service.UserService
}

// Show controlles user information & reservation
func (uc UserController) Show(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	p, err := uc.s.FindUserProfileByID(id.ID)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, p)
}

// PasswordChange controls updating password
func (uc UserController) PasswordChange(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	passwords := entity.UserNewOldPasswords{}
	if err := c.ShouldBindJSON(&passwords); err != nil {
		c.JSON(400, err)
	}

	p, err := uc.s.FindPasswordByID(id.ID)
	if err != nil {
		c.JSON(400, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(p, []byte(passwords.OldPassword)); err != nil {
		c.JSON(400, err)
		return
	}

	if err := uc.s.UpdatePassword(id.ID, passwords.NewPassword); err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.Status(200)
	// TODO: send email
}

// PasswordReset controls resetting password
func (uc UserController) PasswordReset(c *gin.Context) {
	input := entity.UserInputMailPassword{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, err)
		return
	}

	id, err := uc.s.FindIDByEmail(input.Email)
	if err != nil {
		c.JSON(400, err)
		return
	}

	if err := uc.s.UpdatePassword(id, input.Password); err != nil {
		c.AbortWithStatus(404)
	}
	c.Status(200)
	// TODO: send email
}

// UserNameChange chaneg the user name
func (uc UserController) UserNameChange(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	input := entity.UserNameInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, err)
		return
	}

	if err := uc.s.UpdateUserNameByID(id.ID, input.UserName); err != nil {
		c.AbortWithStatus(404)
	}
	c.Status(200)
}

// EmailChange change the user email
func (uc UserController) EmailChange(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
		return
	}

	input := entity.UserEmail{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, err)
		return
	}

	if err := uc.s.UpdateEmailByID(id.ID, input.Email); err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)
	// TODO: send email
}

// // Delete delete the user account
// func (uc UserController) Delete(c *gin.Context) {
// 	id := entity.UserID{}
// 	if err := c.ShouldBindUri(&id); err != nil {
// 		c.JSON(400, err)
// 	}

// 	input := entity.UserEmail{}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(400, err)
// 		return
// 	}

// 	uc.s.FindIDByEmail()

// 	if err := uc.s.DeleteByID(input); err != nil {
// 		c.AbortWithStatus(404)
// 	}
// 	c.Status(200)
// 	// TODO: send email
// }
