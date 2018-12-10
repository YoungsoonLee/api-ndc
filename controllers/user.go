package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/astaxie/beego"

	"github.com/YoungsoonLee/api-ndc/libs"

	"github.com/YoungsoonLee/api-ndc/models"
)

// UserController ...
type UserController struct {
	BaseController
}

type ResetPassword struct {
	ResetToken string `json:"resetToken"`
	Password   string `json:"password"`
}

// ConfirmEmail ...
func (u *UserController) ConfirmEmail() {
	confirmToken := u.GetString(":confirmToken")
	//fmt.Println(confirmToken)

	if len(confirmToken) == 0 {
		u.ResponseError(libs.ErrTokenAbsent, nil)
	}

	// find user by email confirm token
	user, libErr, err := models.CheckConfirmEmailToken(confirmToken)
	if libErr == nil {
		// update
		_, err := models.ConfirmEmail(*user)
		if err != nil {
			u.ResponseError(libs.ErrDatabase, err)
		}
	} else {
		if libErr.Code == "10008" {
			// alaredy confirmed
			u.ResponseSuccess("UID", user.UID)
		} else {
			// error
			u.ResponseError(libErr, err)
		}
	}

	// finish update confirm email.
	// havt to go to login in frontend
	u.ResponseSuccess("UID", user.UID)
}

// ResendConfirmEmail ...
func (u *UserController) ResendConfirmEmail() {
	email := u.GetString(":email")

	// validation
	u.ValidEmail(email)

	// check email
	var user models.User
	user, err := models.FindByEmail(email)
	// if err == nil, already exists Email
	if err != nil {
		u.ResponseError(libs.ErrNoUser, err)
	}

	// update token and send email with confirm token
	user, err = models.ResendConfirmEmail(user)
	if err != nil {
		beego.Error("email confirm update error: ", err)
		u.ResponseError(libs.ErrDatabase, err)
	}

	u.ResponseSuccess("", user)

}

// ForogtPassword ...
func (u *UserController) ForogtPassword() {
	email := u.GetString(":email")

	// validation
	u.ValidEmail(email)

	// check email
	var user models.User
	user, err := models.FindByEmail(email)
	// if err == nil, already exists Email
	if err != nil {
		u.ResponseError(libs.ErrNoUser, err)
	}
	//fmt.Println(user)
	// send forgot password token
	_, err = models.SendPasswordResetToken(user)
	if err != nil {
		u.ResponseError(libs.ErrDatabase, err)
	}

	u.ResponseSuccess("", user)
}

// IsValidResetPasswordToken ...
func (u *UserController) IsValidResetPasswordToken() {
	resetToken := u.GetString(":resetToken")
	//fmt.Println(confirmToken)

	if len(resetToken) == 0 {
		u.ResponseError(libs.ErrTokenAbsent, nil)
	}

	// find user by reset token
	user, libErr, err := models.CheckResetPasswordToken(resetToken)
	if libErr != nil {
		if libErr.Code == "10008" {
			// alaredy confirmed
			u.ResponseSuccess("UID", user.UID)
		} else {
			// error
			u.ResponseError(libErr, err)
		}
	}

	// finish update confirm email.
	// havt to go to login in frontend
	u.ResponseSuccess("UID", user.UID)
}

// ResetPassword ...
func (u *UserController) ResetPassword() {
	var resetPassword ResetPassword
	/*
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &resetPassword)
		if err != nil {
			u.ResponseError(libs.ErrJSONUnmarshal, err)
		}
	*/

	body, _ := ioutil.ReadAll(u.Ctx.Request.Body)
	err := json.Unmarshal(body, &resetPassword)
	if err != nil {
		u.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	if err := models.ResetPassword(resetPassword.ResetToken, resetPassword.Password); err != nil {
		beego.Error("reset password error: ", err)
		u.ResponseError(libs.ErrDatabase, err)
	}

	u.ResponseSuccess("resetToken", resetPassword.ResetToken)
}

// GetProfile ...
func (u *UserController) GetProfile() {
	//var user models.UserFilter
	/*
		UID := u.GetString(":UID")

		beego.Info("UID: ", UID)
		fmt.Println("UID: ", UID)

		// validation
		u.ValidID(UID)

		user, err := models.FindByID(UID)
		if err != nil {
			u.ResponseError(libs.ErrNoUser, err)
		}
		u.ResponseSuccess("", user)
	*/

	//fmt.Println("test")

	et := libs.EasyToken{}
	authtoken := strings.TrimSpace(u.Ctx.Request.Header.Get("Authorization"))
	// new add Bearer
	splitToken := strings.Split(authtoken, "Bearer ")
	if len(splitToken) != 2 {
		u.ResponseError(libs.ErrTokenInvalid, nil)
	}
	valid, uid, err := et.ValidateToken(splitToken[1])

	//beego.Info("Check Login: ", uid, valid)
	//fmt.Println("Check Login: ", uid, valid)

	if !valid || err != nil {
		u.ResponseError(libs.ErrExpiredToken, err)
	}

	// get userinfo
	//var user models.UserFilter
	user, err := models.FindByID(uid)
	if err != nil {
		u.ResponseError(libs.ErrNoUser, err)
	}
	u.ResponseSuccess("", user)
}

// UpdateProfile ...
func (u *UserController) UpdateProfile() {

	et := libs.EasyToken{}
	authtoken := strings.TrimSpace(u.Ctx.Request.Header.Get("Authorization"))

	fmt.Println("authtoken: ", authtoken)

	// new add Bearer
	splitToken := strings.Split(authtoken, "Bearer ")
	fmt.Println("splitToken: ", splitToken, len(splitToken))

	if len(splitToken) != 2 {
		u.ResponseError(libs.ErrTokenInvalid, nil)
	}

	valid, uid, err := et.ValidateToken(splitToken[1])
	if !valid || err != nil {
		u.ResponseError(libs.ErrExpiredToken, err)
	}

	fmt.Println("---1---")

	var user models.User
	user.UID = uid

	fmt.Println("---2---", user)

	err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		u.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	fmt.Println("---3---", user)

	if _, err := models.UpdateProfile(user); err != nil {
		u.ResponseError(libs.ErrDatabase, err)
	}
	u.ResponseSuccess("", user)
}

// UpdatePassword ...
func (u *UserController) UpdatePassword() {
	et := libs.EasyToken{}
	authtoken := strings.TrimSpace(u.Ctx.Request.Header.Get("Authorization"))
	// new add Bearer
	splitToken := strings.Split(authtoken, "Bearer ")
	if len(splitToken) != 2 {
		u.ResponseError(libs.ErrTokenInvalid, nil)
	}
	valid, uid, err := et.ValidateToken(splitToken[1])

	//beego.Info("token: ", authtoken)
	//valid, uid, err := et.ValidateToken(authtoken)

	if !valid || err != nil {
		u.ResponseError(libs.ErrExpiredToken, err)
	}

	var user models.User
	user.UID = uid

	err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		u.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	if _, err := models.UpdatePassword(user); err != nil {
		u.ResponseError(libs.ErrDatabase, err)
	}
	u.ResponseSuccess("", user)

}

// ---------------------------------------------------------------------------------------------------------------
