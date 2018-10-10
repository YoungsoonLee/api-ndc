package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
)

type BillingController struct {
	BaseController
}

// Xsolla struct
type XSuser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Ip    string `json:"ip"`
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
	err := json.Unmarshal(b.Ctx.Input.RequestBody, &pt)
	if err != nil {
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// validation param

	pt, err = models.AddPaymentTry(pt)
	if err != nil {
		b.ResponseError(libs.ErrDatabase, err)
	}

	//fmt.Println(pt)

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

	signature := b.Ctx.Request.Header.Get("Authorization Signature ")
	xsollaData.Signature = signature

	body, _ := ioutil.ReadAll(b.Ctx.Request.Body)
	err := json.Unmarshal(body, &xsollaData)
	if err != nil {
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	fmt.Println("xsollaData.Signature: ", xsollaData.Signature)
	fmt.Println("xsollaData.NotificationType: ", xsollaData.NotificationType)
	fmt.Println("xsollaData.Purchase.Total.Amount: ", xsollaData.Purchase.Total.Amount)
	fmt.Println("xsollaData.Purchase.Total.Currency: ", xsollaData.Purchase.Total.Currency)
	fmt.Println("xsollaData.Signature: ", xsollaData.Signature)
	fmt.Println("xsollaData.Transaction.ExternalID: ", xsollaData.Transaction.ExternalID)
	fmt.Println("xsollaData.Transaction.ID: ", xsollaData.Transaction.ID)
	fmt.Println("xsollaData.Transaction.PaymentDate: ", xsollaData.Transaction.PaymentDate)
	fmt.Println("xsollaData.User.ID: ", xsollaData.User.ID)
	fmt.Println("xsollaData.User.Email: ", xsollaData.User.Email)
	fmt.Println("xsollaData.User.Ip: ", xsollaData.User.Ip)

	b.ResponseSuccess("", "")

}
