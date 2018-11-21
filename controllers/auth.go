package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
	"github.com/astaxie/beego"
)

// AuthController ...
type AuthController struct {
	BaseController
}

// LoginToken ...
type LoginToken struct {
	Displayname string `json:"displayname"`
	UID         int64  `json:"uid"`
	Token       string `json:"token"`
}

// Social ...
type Social struct {
	Provider            string `json:"provider"`
	ProviderAccessToken string `json:"accessToken"`
	Email               string `json:"email"`
	ProviderID          string `json:"providerId"`
	Picture             string `json:"picture"`
}

// AuthedData ...
type AuthedData struct {
	UID         int64  `json:"uid"`
	Displayname string `json:"displayname"`
	Balance     int    `json:"balance"`
	Picture     string `json:"picture"`
}

// CheckDisplayName ...
// @Title CheckDisplayName
// @Description create services
// @Param	body		body 		true		""
// @Success 200 {string} displayname
// @Failure 403 body is empty
// @router /:displayname [get]
func (c *AuthController) CheckDisplayName() {

	displayname := c.GetString(":displayname")
	// validation
	c.ValidDisplayname(displayname)

	_, err := models.FindByDisplayname(displayname)
	// if err == nil, already exists displayname
	if err == nil {
		c.ResponseError(libs.ErrDupDisplayname, err)
	}

	//success
	c.ResponseSuccess("displayname", displayname)
}

// CreateUser ...
// @Title CreateUser except social
// @Description create users
// @Param	displayname		query 	string	true		"displayname"
// @Param	email			query 	string	true		"email"
// @Param	password		query 	string	true		"password"
// @Success 200 {int} models.User.UID
// @Failure 403 body is empty
// @router /CreateUser [post]
func (c *AuthController) CreateUser() {
	var user models.User
	/*
		err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
		if err != nil {
			c.ResponseError(libs.ErrJSONUnmarshal, err)
		}
	*/
	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		c.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// validation
	c.ValidDisplayname(user.Displayname)
	c.ValidEmail(user.Email)
	c.ValidPassword(user.Password)

	// check dup displayname
	_, err = models.FindByDisplayname(user.Displayname)

	// if err == nil, already exists displayname
	if err == nil {
		c.ResponseError(libs.ErrDupDisplayname, err)
	}

	// check dup email
	_, err = models.FindByEmail(user.Email)
	// if err == nil, already exists Email
	if err == nil {
		c.ResponseError(libs.ErrDupEmail, err)
	}

	// save to db
	UID, err := models.AddUser(user)
	if err != nil {
		c.ResponseError(libs.ErrDatabase, err)
	}

	fmt.Println("CreateUser: ", UID)

	// auto login
	user.UID = UID
	c.makeLogin(&user)
}

// Login ...
// @Title Login
// @Description Logs user into the system
// @Param	displayname		query 	string	true		"displayname for login"
// @Param	password		query 	string	true		"password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (c *AuthController) Login() {
	var user models.User

	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		c.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// validation
	inputPass := user.Password
	c.ValidDisplayname(user.Displayname)
	c.ValidPassword(user.Password)

	// Find salt, password hash for auth
	user, err = models.FindAuthByDisplayname(user.Displayname)
	if err != nil {
		c.ResponseError(libs.ErrPass, err)
	}

	beego.Info(user)

	if user.Provider == "facebook" && user.Password == "" {
		c.ResponseError(libs.ErrLoginFacebook, nil)
	}
	if user.Provider == "google" && user.Password == "" {
		c.ResponseError(libs.ErrLoginGoogle, nil)
	}

	// check password
	ok, err := user.CheckPass(inputPass)
	if !ok || err != nil {
		// wrong password
		c.ResponseError(libs.ErrPass, err)
	}

	//beego.Info(user)
	c.makeLogin(&user)
}

// CheckLogin ...
func (c *AuthController) CheckLogin() {

	et := libs.EasyToken{}
	authtoken := strings.TrimSpace(c.Ctx.Request.Header.Get("Authorization"))
	valid, uid, err := et.ValidateToken(authtoken)

	beego.Info("Check Login: ", uid, valid)
	fmt.Println("Check Login: ", uid, valid)

	if !valid || err != nil {
		c.ResponseError(libs.ErrExpiredToken, err)
	}

	// get userinfo
	//var user models.UserFilter
	user, err := models.FindByID(uid)
	if err != nil {
		c.ResponseError(libs.ErrNoUser, err)
	}

	//beego.Info(user)

	c.ResponseSuccess("", AuthedData{user.UID, user.Displayname, user.Balance, user.Picture})
}

// Social ...
// @Title CreateUser or SigninUser for social FB and G+
// @Description create social users or signin
// @Param	provider		query 	string	true		"The provider (FB, G+)"
// @Param	accessToken		query 	string	true		"The accessToken"
// @Success 200 {int}
// @Failure 403 body is empty
// @router /Social [post]
func (c *AuthController) Social() {
	var social Social
	/*
		err := json.Unmarshal(c.Ctx.Input.RequestBody, &social)
		if err != nil {
			c.ResponseError(libs.ErrJSONUnmarshal, err)
		}
	*/
	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	err := json.Unmarshal(body, &social)
	if err != nil {
		c.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// TODO: validation
	// unless provier is null or accessToken is null, get error

	var user models.User
	user, err = models.FindByEmail(social.Email)

	// if err == nil, already exists Email
	if err == nil {
		// make login
		//fmt.Println("already exists email ", user)
		//update social info, it can login local and social both.
		if len(user.Provider) == 0 || user.Provider != social.Provider {
			user.Provider = social.Provider
			user.ProviderAccessToken = social.ProviderAccessToken
			user.ProviderID = social.ProviderID
			user.Picture = social.Picture
			user.Confirmed = true

			c.updateSocialInfo(user)
		}

		c.makeLogin(&user)

	} else {
		// add social user
		user.Provider = social.Provider
		user.ProviderAccessToken = social.ProviderAccessToken
		user.ProviderID = social.ProviderID
		user.Email = social.Email
		user.Picture = social.Picture
		c.createSocialUser(user)
	}

}

// @Title Logout
// @Description Logs user into the system
// @Param	displayname		query 	string	true		"The displayname for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (c *AuthController) Logout() {

}

func (c *AuthController) createSocialUser(user models.User) {

	UID, displayname, err := models.AddSocialUser(user)
	if err != nil {
		c.ResponseError(libs.ErrDatabase, err)
	}

	user.UID = UID
	user.Displayname = displayname
	c.makeLogin(&user)
}

func (c *AuthController) updateSocialInfo(user models.User) {
	UID, displayname, err := models.UpdateSocialInfo(user)
	if err != nil {
		c.ResponseError(libs.ErrDatabase, err)
	}

	user.UID = UID
	user.Displayname = displayname
	c.makeLogin(&user)
}

func (c *AuthController) makeLogin(user *models.User) {
	fmt.Println("makeLogin: ", user.UID)

	// login
	et := libs.EasyToken{
		Displayname: user.Displayname,
		UID:         user.UID,
		Expires:     time.Now().Unix() + 3600*24*7, // 7days, 1 hour(3600)
	}

	token, err := et.GetToken()
	if token == "" || err != nil {
		c.ResponseError(libs.ErrTokenOther, nil)
	}

	//beego.Info("makeLogin: ", user.UID)
	c.ResponseSuccess("", LoginToken{user.Displayname, user.UID, token})
}
