package controllers

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
)

type BillingController struct {
	BaseController
}

// Xsolla struct
type XSuser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Ip      string `json:"ip"`
	Phone   string `json:"phone"`
	Country string `json:"country"`
}

type XSpurchaseDetail struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}

type XSpurchase struct {
	Total XSpurchaseDetail
}

type XStransaction struct {
	ID          int       `json:"id"`           // TxID from Xsolla
	ExternalID  string    `json:"external_id"`  // PxID
	PaymentDate time.Time `json:"payment_date"` // transaction_at
}

//  xsolla callback data
type XSollaData struct {
	Signature        string        `json:"signature"`
	NotificationType string        `json:"notification_type"`
	User             XSuser        `json:"user"`
	Purchase         XSpurchase    `json:"purchase"`
	Transaction      XStransaction `json:"transaction"`
}

// GetChargeItems ...
// @Title Create Payment Category
// @Description create payment category
// @Success 200 {int} models.PaymentItem.ItemID
// @Failure 403 body is empty
// @router / [GET]
func (b *BillingController) GetChargeItems() {

	// save to db
	chargeItems, err := models.GetChargeItems()
	if err != nil {
		b.ResponseError(libs.ErrDatabase, err)
	}

	//success
	b.ResponseSuccess("", chargeItems)
}

// GetPaymentToken ...
// @Title Get PaymentToken
// @Description create payment category
// @Param	UID			json 	INT		false		"user id"
// @Param	ItemID		json 	INT		false		"item id"
// @Success 200 {int} models.PaymentItem.PayTryID
// @Failure 403 body is empty
// @router / [post]
func (b *BillingController) GetPaymentToken() {
	//
	var pt models.PaymentTry

	body, _ := ioutil.ReadAll(b.Ctx.Request.Body)
	//err := json.Unmarshal(b.Ctx.Input.RequestBody, &pt)
	err := json.Unmarshal(body, &pt)
	if err != nil {
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// validation param

	// insert payment try
	pt, err = models.AddPaymentTry(pt)
	if err != nil {
		b.ResponseError(libs.ErrDatabase, err)
	}

	url := os.Getenv("XSOLLA_ENDPOINT") + os.Getenv("XSOLLA_MERCHANT_ID") + "/token"
	beego.Info("url: ", url)

	// make json send data for getting token
	var sendDataToGetToken libs.XsollaSendJSONToGetToken
	sendDataToGetToken.User.ID.Value = pt.UID
	sendDataToGetToken.User.ID.Hidden = true
	sendDataToGetToken.User.Email.Value = "" // TODO: ???
	sendDataToGetToken.User.Email.AllowModify = false
	sendDataToGetToken.User.Email.Hidden = true
	sendDataToGetToken.User.Country.Value = "US"
	sendDataToGetToken.User.Name.Value = "" // TODO: ??? nickname
	sendDataToGetToken.User.Name.Hidden = false

	sendDataToGetToken.Settings.ProjectID = 24380
	sendDataToGetToken.Settings.ExternalID = pt.PxID
	sendDataToGetToken.Settings.Mode = pt.Mode
	sendDataToGetToken.Settings.Language = "en"
	sendDataToGetToken.Settings.Currency = "USD"
	sendDataToGetToken.Settings.UI.Size = "medium"

	sendDataToGetToken.Purchase.Checkout.Currency = "USD"
	sendDataToGetToken.Purchase.Checkout.Amount = pt.Price // price
	sendDataToGetToken.Purchase.Description.Value = pt.ItemName

	sendDataToGetToken.CustomParameters.Pid = pt.PxID

	jsonStr, err := json.Marshal(sendDataToGetToken)
	if err != nil {
		beego.Error("marshall error: ", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		beego.Error("NewRequest error: ", err)
	}

	req.Header.Set("Content-Type", "application.json")

	setHeaderKey := "Basic " + os.Getenv("XSOLLA_API_KEY")
	beego.Info("setHeaderKey: ", setHeaderKey)

	req.Header.Set("Authorization", setHeaderKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("client error: ", err)
	}

	fmt.Println("response status: ", resp.Status)
	fmt.Println("response header: ", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ", body)

	b.ResponseSuccess("", pt)

}

// CallbackXsolla ...
// @Title Get xsolla callback data
// @Description xsolla send callbac data
// ...
// ...
// ...
// ...
func (b *BillingController) CallbackXsolla() {
	var xsollaData XSollaData

	signature := strings.TrimSpace(b.Ctx.Request.Header.Get("Authorization"))
	signature = strings.Replace(signature, "Signature ", "", -1)
	if signature == "" {
		b.XsollaResponseError(libs.ErrXNilSig)
	}

	xsollaData.Signature = signature

	body, _ := ioutil.ReadAll(b.Ctx.Request.Body)
	if body == nil {
		body = b.Ctx.Input.RequestBody // for local test
	}

	err := json.Unmarshal(body, &xsollaData)
	if err != nil {
		//b.ResponseError(libs.ErrJSONUnmarshal, err)
		b.XsollaResponseError(libs.ErrXInvalidJSON)
	}

	beego.Info("xsollaData: ", xsollaData)

	// hashed
	h := sha1.New()
	hBody := string(body) + os.Getenv("XSOLLA_SECRET_KEY")
	h.Write([]byte(hBody))
	hashedData := fmt.Sprintf("%x", h.Sum(nil))

	if hashedData != xsollaData.Signature {
		beego.Error(hashedData, xsollaData.Signature)
		b.XsollaResponseError(libs.ErrXInvalidSig)
	}

	// check user
	_, err = models.FindByID(xsollaData.User.ID)
	if err != nil {
		b.XsollaResponseError(libs.ErrXInvalidUser)
	}

	// check notification_type == "user_validation"
	if xsollaData.NotificationType == "user_validation" {
		b.ResponseSuccess("", "") //success
	}

	// check notification_type == "payment"
	if xsollaData.NotificationType == "payment" {
		// check payment try
		var pt models.PaymentTry
		pt.PxID = xsollaData.Transaction.ExternalID
		pt.UID, _ = strconv.ParseInt(xsollaData.User.ID, 10, 64)
		pt.Amount = xsollaData.Purchase.Total.Amount
		pt, exists := models.CheckPaymentTry(pt)
		if !exists {
			b.XsollaResponseError(libs.ErrXInvalidPaytryData)
		}

		// make charge data
		var c models.PaymentTransaction
		c.PxID = xsollaData.Transaction.ExternalID
		c.TxID = strconv.Itoa(xsollaData.Transaction.ID)
		c.UID, _ = strconv.ParseInt(xsollaData.User.ID, 10, 64)
		c.ItemID = pt.ItemID
		c.ItemName = pt.ItemName
		c.PgID = pt.PgID
		c.Currency = pt.Currency
		c.Price = pt.Price
		c.Amount = pt.Amount

		beego.Info("charge data: ", c)

		// begin tran
		err := models.AddPaymentTransaction(c)
		if err != nil {
			beego.Error("Charge error: ", err)
			b.XsollaResponseError(libs.ErrXMakePaytransaction)
		}

		// set redis?
		// use ??
		//

		//success
		b.ResponseSuccess("", "")

	} else {
		// invalid paytry data
		b.XsollaResponseError(libs.ErrXInvalidNotiType)
	}

	/*
		fmt.Println("xsollaData.Signature: ", xsollaData.Signature)
		fmt.Println("xsollaData.NotificationType: ", xsollaData.NotificationType)
		fmt.Println("xsollaData.Purchase.Total.Amount: ", xsollaData.Purchase.Total.Amount)
		fmt.Println("xsollaData.Purchase.Total.Currency: ", xsollaData.Purchase.Total.Currency)
		fmt.Println("xsollaData.Transaction.ExternalID: ", xsollaData.Transaction.ExternalID)
		fmt.Println("xsollaData.Transaction.ID: ", xsollaData.Transaction.ID)
		fmt.Println("xsollaData.Transaction.PaymentDate: ", xsollaData.Transaction.PaymentDate)
		fmt.Println("xsollaData.User.ID: ", xsollaData.User.ID)
		fmt.Println("xsollaData.User.Email: ", xsollaData.User.Email)
		fmt.Println("xsollaData.User.Ip: ", xsollaData.User.Ip)
		fmt.Println("xsollaData.User.Phone: ", xsollaData.User.Phone)
		fmt.Println("xsollaData.User.Country: ", xsollaData.User.Country)
	*/
}
