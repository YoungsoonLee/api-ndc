package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/zoonman/gravatar"
	"golang.org/x/crypto/scrypt"
)

// User ...
type User struct {
	UID                 int64      `orm:"column(UID);pk" json:"uid"`          // user id
	Displayname         string     `orm:"size(30);unique" json:"displayname"` // 4 ~ 16 letters for local,
	Email               string     `orm:"size(100);unique" json:"email"`      // max 100 letters
	Password            string     `orm:"null" json:"password"`               // if account is provider, this column is null
	Salt                string     `orm:"null" json:"salt"`
	PasswordResetToken  string     `orm:"size(1000);null" json:"password_reset_token"`
	PasswordResetExpire *time.Time `orm:"null"  json:"password_reset_expire"`
	Confirmed           bool       `orm:"default(false)" json:"confirmed"`
	ConfirmResetToken   string     `orm:"size(1000);null" json:"confirm_reset_token"`
	ConfirmResetExpire  time.Time  `orm:"null"  json:"confirm_reset_expire"`
	Picture             string     `orm:"size(1000);null" json:"picture"`
	Provider            string     `orm:"size(50);null" json:"provider"` // google , facebook
	ProviderID          string     `orm:"column(ProviderID);size(1000);null" json:"provider_id"`
	ProviderAccessToken string     `orm:"size(1000);null" json:"provider_access_token"`
	Permission          string     `orm:"size(50);default(user)" json:"permission"`     // user, admin ...
	Status              string     `orm:"size(50);default(normal)" json:"status"`       // normal, ban, close ...
	CreateAt            time.Time  `orm:"type(datetime);auto_now_add" json:"create_at"` // first save
	UpdateAt            time.Time  `orm:"type(datetime);auto_now" json:"update_at"`     // eveytime save
	Balance             int        `orm:"-" json:"balance"`                             // wallet's balance
}

// UserFilter ...
// for giving user's info to front or game
type UserFilter struct {
	UID         int64     `orm:"column(UID)" json:"uid"`       // user id
	Displayname string    `orm:"size(30);" json:"displayname"` // 4 ~ 16 letters for local,
	Email       string    `orm:"size(100);" json:"email"`      // max 100 letters
	Picture     string    `orm:"size(1000);" json:"picture"`
	Provider    string    `orm:"size(50);" json:"provider"`        // google , facebook
	Permission  string    `orm:"size(50);" json:"permission"`      // user, admin ...
	Status      string    `orm:"size(50);" json:"status"`          // normal, ban, close ...
	CreateAt    time.Time `orm:"type(datetime);" json:"create_at"` // first save
	UpdateAt    time.Time `orm:"type(datetime);" json:"update_at"` // eveytime save
	Balance     int       `orm:"-" json:"balance"`                 // wallet's balance
}

const pwHashBytes = 64

func generateSalt() (salt string, err error) {
	buf := make([]byte, pwHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

func generatePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, pwHashBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h), nil
}

// CheckPass compare input password.
func (u *User) CheckPass(pass string) (bool, error) {
	//fmt.Println(pass, u.Salt)
	hash, err := generatePassHash(pass, u.Salt)
	if err != nil {
		return false, err
	}

	return u.Password == hash, nil
}

// AddUser ...
func AddUser(u User) (int64, error) {

	// make Id
	u.UID = time.Now().UnixNano()

	// make hashed password
	salt, err := generateSalt()
	if err != nil {
		return 0, err
	}
	hash, err := generatePassHash(u.Password, salt)
	if err != nil {
		return 0, err
	}

	// set password & salt
	u.Password = hash
	u.Salt = salt

	// get gravatar
	u.Picture = gravatar.Avatar(u.Email, 80)

	// make email confirm token
	u2, err := uuid.NewV4()
	if err != nil {
		return 0, err
	}
	u.ConfirmResetToken = u2.String()
	u.ConfirmResetExpire = time.Now().Add(1 * time.Hour)

	fmt.Println("insert add user: ", u.UID)

	// save to db with transaction user and wallet
	o := orm.NewOrm()
	err = o.Begin()

	_, err = o.Insert(&u)
	if err != nil {
		err = o.Rollback()
		return 0, err
	}

	wallet := Wallet{UID: u.UID, Balance: 0}
	_, err = o.Insert(&wallet)
	if err != nil {
		err = o.Rollback()
		return 0, err
	}

	err = o.Commit()

	// send confirm mail async
	go libs.MakeMail(u.Email, "confirm", u.ConfirmResetToken)

	fmt.Println("before return add user: ", u.UID)
	return u.UID, nil
}

// AddSocialUser ...
func AddSocialUser(u User) (int64, string, error) {
	// make Id
	u.UID = time.Now().UnixNano()

	// for displayname
	b := make([]byte, 5) //equals 8 charachters
	rand.Read(b)
	s := hex.EncodeToString(b)
	u.Displayname = s

	u.Confirmed = true

	// save to db with transaction user and wallet
	o := orm.NewOrm()
	err := o.Begin()

	_, err = o.Insert(&u)
	if err != nil {
		err = o.Rollback()
		return 0, "", err
	}

	wallet := Wallet{UID: u.UID, Balance: 0}
	_, err = o.Insert(&wallet)
	if err != nil {
		err = o.Rollback()
		return 0, "", err
	}

	err = o.Commit()

	return u.UID, u.Displayname, nil
}

// UpdateSocialInfo ...
func UpdateSocialInfo(u User) (int64, string, error) {

	o := orm.NewOrm()
	if _, err := o.Update(&u, "Provider", "ProviderAccessToken", "ProviderID", "Picture", "Confirmed"); err != nil {

		return 0, "", err
	}

	return u.UID, u.Displayname, nil
}

// FindAuthByDisplayname ...
// using for auth
func FindAuthByDisplayname(displayname string) (User, error) {
	var user User
	o := orm.NewOrm()

	sql := "SELECT " +
		" \"UID\" , " +
		" Displayname, " +
		" Password, " +
		" Salt, " +
		" Provider " +
		" FROM \"user\" " +
		" WHERE Displayname = ?"
	err := o.Raw(sql, displayname).QueryRow(&user)
	return user, err
}

// FindByDisplayname ...
// TODO: add balance
func FindByDisplayname(displayname string) (User, error) {
	var user User
	o := orm.NewOrm()

	sql := "SELECT " +
		" \"UID\" , " +
		" Displayname, " +
		" Email, " +
		" Confirmed, " +
		" Picture, " +
		" Provider, " +
		" Permission, " +
		" Status, " +
		" Create_At, " +
		" Update_At " +
		" FROM \"user\" " +
		" WHERE Displayname = ?"
	err := o.Raw(sql, displayname).QueryRow(&user)

	return user, err
}

// FindByEmail ...
// TODO: add balance
func FindByEmail(email string) (User, error) {
	var user User
	o := orm.NewOrm()

	sql := "SELECT " +
		" \"UID\" , " +
		" Displayname, " +
		" Email, " +
		" Confirmed, " +
		" Picture, " +
		" Provider, " +
		" Permission, " +
		" Status, " +
		" Create_At, " +
		" Update_At " +
		" FROM \"user\" " +
		" WHERE Email = ?"
	err := o.Raw(sql, email).QueryRow(&user)
	return user, err
}

// FindByID ...
func FindByID(id string) (*UserFilter, error) {
	beego.Info("findByID: ", id)
	var user *UserFilter

	o := orm.NewOrm()
	sql := "SELECT " +
		" \"user\".\"UID\" , " +
		" Displayname, " +
		" Email, " +
		" Confirmed, " +
		" Picture, " +
		" Provider, " +
		" Permission, " +
		" Status, " +
		" \"user\".Create_At, " +
		" \"user\".Update_At, " +
		" \"wallet\".Balance " +
		" FROM \"user\", \"wallet\" " +
		" WHERE \"user\".\"UID\" = \"wallet\".\"UID\" " +
		" and \"user\".\"UID\" = ? "

	err := o.Raw(sql, id).QueryRow(&user)
	return user, err
}

// FindByProvider ...
func FindByProvider(provider string, accessToken string, providerID string) bool {

	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("Provider", provider).Filter("ProviderAccessToken", accessToken).Filter("ProviderID", providerID).Exist()

	return exist
}

// CheckConfirmEmailToken ...
func CheckConfirmEmailToken(token string) (*User, *libs.ControllerError, error) {
	var user *User
	o := orm.NewOrm()

	sql := "SELECT " +
		" \"UID\" , " +
		" Displayname, " +
		" Confirmed " +
		" FROM \"user\" " +
		" WHERE Confirm_Reset_Token = ? AND Confirmed = true"
	err := o.Raw(sql, token).QueryRow(&user)

	if err == nil {
		// already confirmed or wrong token
		return user, libs.ErrAlreadyConfirmed, err
	}

	// wrong token
	err = o.Raw("select \"UID\", Displayname, Confirmed from \"user\" where Confirm_Reset_Token =? and Confirmed = false", token).QueryRow(&user)
	if err != nil {
		// already confirmed or wrong token
		return user, libs.ErrWrongToken, err
	}

	//  expired token
	err = o.Raw("select \"UID\", Displayname, Confirmed from \"user\" where Confirm_Reset_Token =? and Confirm_Reset_Expire <= ?", token, time.Now()).QueryRow(&user)
	if err == nil {
		// expire token
		return user, libs.ErrExpiredToken, err
	}

	return user, nil, nil
}

// Confirm Email...
// func ConfirmEmail(token string) (User, error) {
func ConfirmEmail(u User) (User, error) {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE \"user\" SET Confirmed = ?, Confirm_Reset_Expire =? WHERE \"UID\"=?", true, nil, u.UID).Exec()
	if err != nil {
		return User{}, err
	}

	return u, err
}

func ResendConfirmEmail(u User) (User, error) {
	// make email confirm token
	u2, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}

	u.ConfirmResetToken = u2.String()
	u.ConfirmResetExpire = time.Now().Add(1 * time.Hour)
	u.Confirmed = false

	o := orm.NewOrm()
	if _, err := o.Update(&u, "Confirmed", "ConfirmResetToken", "ConfirmResetExpire"); err != nil {
		return User{}, err
	}

	// send confirm mail async
	go libs.MakeMail(u.Email, "confirm", u.ConfirmResetToken)

	return u, nil
}

func SendPasswordResetToken(u User) (User, error) {
	// make forgot password token
	u2, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}

	u.PasswordResetToken = u2.String()
	ct := time.Now().Add(1 * time.Hour)
	u.PasswordResetExpire = &ct

	o := orm.NewOrm()
	if _, err := o.Update(&u, "PasswordResetToken", "PasswordResetExpire"); err != nil {
		return User{}, err
	}

	// send confirm mail async
	go libs.MakeMail(u.Email, "forgotPassword", u.PasswordResetToken)

	return u, nil
}

func CheckResetPasswordToken(resetToken string) (*User, *libs.ControllerError, error) {
	var user *User

	o := orm.NewOrm()
	// wrong token
	err := o.Raw("select \"UID\", Displayname, Confirmed from \"user\" where Password_Reset_Token =?", resetToken).QueryRow(&user)
	if err != nil {
		// already confirmed or wrong token
		return user, libs.ErrTokenInvalid, err
	}

	//  expired token
	err = o.Raw("select \"UID\", Displayname, Confirmed from \"user\" where Password_Reset_Token =? and Password_Reset_Expire <= ?", resetToken, time.Now()).QueryRow(&user)
	if err == nil {
		// expire tokens
		return user, libs.ErrExpiredToken, err
	}

	return user, nil, nil
}

// ResetPassword ...
func ResetPassword(resetToken, password string) error {

	// make hashed password
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	hash, err := generatePassHash(password, salt)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	_, err = o.Raw("UPDATE \"user\" SET Password = ?, Salt = ?, Password_Reset_Token = ?, Password_Reset_Expire=? WHERE Password_Reset_Token = ?", hash, salt, nil, nil, resetToken).Exec()
	//fmt.Println(r.LastInsertId())
	if err != nil {
		return err
	}

	return nil
}

// UpdateProfile
func UpdateProfile(u User) (User, error) {
	//TODO: if email changed, send email confirm.

	o := orm.NewOrm()
	if _, err := o.Update(&u, "Displayname", "Email"); err != nil {
		return User{}, err
	}

	return u, nil
}

// UpdatePassword
func UpdatePassword(u User) (User, error) {
	o := orm.NewOrm()

	// make hashed password
	salt, err := generateSalt()
	if err != nil {
		return User{}, err
	}
	hash, err := generatePassHash(u.Password, salt)
	if err != nil {
		return User{}, err
	}

	// set password & salt
	u.Password = hash
	u.Salt = salt

	if _, err := o.Update(&u, "Password", "Salt"); err != nil {
		return User{}, err
	}

	return u, nil
}

// ---------------------------------------------------------------------------------------------------------------
