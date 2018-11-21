package controllers

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
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

// xsolla callback data
type XSollaData struct {
	Signature        string        `json:"signature"`
	NotificationType string        `json:"notification_type"`
	User             XSuser        `json:"user"`
	Purchase         XSpurchase    `json:"purchase"`
	Transaction      XStransaction `json:"transaction"`
}

// deduct
type DeductInput struct {
	/***
	 * Inputs ...
	 * 	service_id: 각 게임 별 할당 되는 고유 ID
	 * 	access_toke: 유저 ID
	 * 	external_txid: 각 게임 서비스 고유의 트랜잭션 ID
	 * 	external_itemid: 각 게임 서비스의 구매시의 해당 아이템 ID. (조회, 통계, 추적용)
	 * 	external_itemname: 각 게임 서비스의 구매시의 해당 아이템 이름. (조회, 통계, 추적용)
	 * 	item_amount: 차감 해야 될 cyber coin 양
	 *
	 * 	요청시 헤더 Authorization: Signature에 sha1(위 input을 json으로 + 배포되는 service_key)
	 *
	 * TODO: user's ip ???
	 *
	 * Outputs...
	 * 	service_id: 각 게임 별 할당 되는 고유 ID
	 * 	external_id: 각 게임 서비스 고유의 트랜잭션 ID
	 * 	deduct_id: cyber coin 차감 후 발생한 고유 트랜잭션 ID
	 *
	 */
	ServiceID        string `json:"service_id"`
	ExternalTxID     string `json:"external_txid"`
	ExternalItemID   string `json:"external_itemid"`
	ExternalItemName string `json:"external_itemname"`
	Amount           int    `json:"amount"`
	Hash             string `json:"hash"`
}

type HashedBody struct {
	ServiceID        string `json:"service_id"`
	ExternalTxID     string `json:"external_txid"`
	ExternalItemID   string `json:"external_itemid"`
	ExternalItemName string `json:"external_itemname"`
	Amount           int    `json:"amount"`
}

type ResultDeduct struct {
	UID              int64  `json:"uid"`
	ServiceID        string `json:"service_id"`
	ExternalTxID     string `json:"external_txid"`
	ExternalItemID   string `json:"external_itemid"`
	ExternalItemName string `json:"external_itemname"`
	Amount           int    `json:"amount"`
	BalanceAfterBuy  int    `json:"balance_after_buy"`
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
	err := json.Unmarshal(body, &pt)
	if err != nil {
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// validation param uid
	// check UID
	//var user models.UserFilter
	user, err := models.FindByID(strconv.FormatInt(pt.UID, 10))
	if err != nil {
		b.ResponseError(libs.ErrNoUser, err)
	}

	// insert payment try
	pt, err = models.AddPaymentTry(pt)
	if err != nil {
		b.ResponseError(libs.ErrDatabase, err)
	}

	url := os.Getenv("XSOLLA_ENDPOINT") + os.Getenv("XSOLLA_MERCHANT_ID") + "/token"
	// beego.Info("url: ", url)

	// make json send data for getting token
	var sendDataToGetToken libs.XsollaSendJSONToGetToken
	sendDataToGetToken.User.ID.Value = strconv.FormatInt(pt.UID, 10)
	sendDataToGetToken.User.ID.Hidden = true
	sendDataToGetToken.User.Email.Value = user.Email
	sendDataToGetToken.User.Email.AllowModify = false
	sendDataToGetToken.User.Email.Hidden = true
	sendDataToGetToken.User.Country.Value = "US"
	sendDataToGetToken.User.Name.Value = user.Displayname
	sendDataToGetToken.User.Name.Hidden = false

	sendDataToGetToken.Settings.ProjectID = 24380
	sendDataToGetToken.Settings.ExternalID = pt.PxID
	sendDataToGetToken.Settings.Mode = pt.Mode
	sendDataToGetToken.Settings.Language = "en"
	sendDataToGetToken.Settings.Currency = "USD"
	sendDataToGetToken.Settings.UI.Size = "medium"

	sendDataToGetToken.Purchase.Checkout.Currency = "USD"
	sendDataToGetToken.Purchase.Checkout.Amount = float32(pt.Price) // price
	sendDataToGetToken.Purchase.Description.Value = pt.ItemName

	sendDataToGetToken.CustomParameters.Pid = pt.PxID

	jsonStr, err := json.Marshal(sendDataToGetToken)
	if err != nil {
		beego.Error("sendDataToGetToken marshall error: ", err)
		b.ResponseError(libs.ErrJSONmarshal, err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		beego.Error("getToekn Request error: ", err)
		b.ResponseError(libs.ErrTokenRequest, err)
	}

	key := os.Getenv("XSOLLA_MERCHANT_ID") + ":" + os.Getenv("XSOLLA_API_KEY")
	encoded := base64.StdEncoding.EncodeToString([]byte(key))
	setHeaderKey := "Basic " + encoded
	// beego.Info("setHeaderKey: ", setHeaderKey, os.Getenv("XSOLLA_API_KEY"))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", setHeaderKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("client error: ", err)
		b.ResponseError(libs.ErrClient, err)
	}

	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &pt)
	if err != nil {
		beego.Error("get token unmarshall error: ", err)
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	beego.Info("token: ", pt.Token)

	// TODO: check token is nil

	b.ResponseSuccess("", pt)

}

// GetChargeHistory ..
func (b *BillingController) GetChargeHistory() {
	UID := b.GetString(":UID")
	// TODO: validation.
	if UID == "" {
		err := errors.New("UID is nil")
		b.ResponseError(libs.ErrInputData, err)
	}
	iUID, _ := strconv.ParseInt(UID, 10, 64)
	paytransacsion, err := models.GetPayTransaction(iUID)
	if err != nil {
		b.ResponseError(libs.ErrDatabase, err)
	}

	//fmt.Println(paytransacsion)

	b.ResponseSuccess("tabulator", paytransacsion)
}

// GetUsedHistory ...
func (b *BillingController) GetUsedHistory() {
	UID := b.GetString(":UID")
	// TODO: validation.
	if UID == "" {
		err := errors.New("UID is nil")
		b.ResponseError(libs.ErrInputData, err)
	}
	iUID, _ := strconv.ParseInt(UID, 10, 64)
	duductHistory, err := models.GetUsedHistory(iUID)
	if err != nil {
		b.ResponseError(libs.ErrDatabase, err)
	}

	//fmt.Println(duductHistory)

	b.ResponseSuccess("tabulator", duductHistory)
}

// BuyItem ...
// deduct cyber coin
func (b *BillingController) BuyItem() {
	/***
	 * Inputs ...
	 * 	service_id: 각 게임 별 할당 되는 고유 ID
	 * 	external_id: 각 게임 서비스 고유의 트랜잭션 ID
	 * 	item_name: 각 게임 서비스의 구매시의 해당 아이템 이름. (조회, 통계, 추적용)
	 * 	item_id: 각 게임 서비스의 구매시의 해당 아이템 ID. (조회, 통계, 추적용)
	 * 	item_amount: 차감 해야 될 cyber coin 양
	 * 	hash: sha1(위 input을 json으로 + 배포되는 service_key)
	 *
	 * 	요청시 헤더 Authorization: access_token
	 *
	 * TODO: user's ip ???
	 *
	 * Outputs...
	 * 	service_id: 각 게임 별 할당 되는 고유 ID
	 * 	external_id: 각 게임 서비스 고유의 트랜잭션 ID
	 * 	deduct_id: cyber coin 차감 후 발생한 고유 트랜잭션 ID
	 */

	start := time.Now()

	// TODO: need more performance !!!
	// get header for auth
	authtoken := strings.TrimSpace(b.Ctx.Request.Header.Get("Authorization"))
	if authtoken == "" {
		b.ResponseError(libs.ErrTokenAbsent, errors.New(libs.ErrTokenAbsent.Message))
	}
	// check UID, check valid token
	et := libs.EasyToken{}
	valid, uid, err := et.ValidateToken(authtoken)
	if !valid || err != nil {
		b.ResponseError(libs.ErrExpiredToken, err)
	}

	// get body
	var deductInput DeductInput
	body, _ := ioutil.ReadAll(b.Ctx.Request.Body)
	err = json.Unmarshal(body, &deductInput)
	if err != nil {
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// TODO: make log file for inputs... with go routine ??s
	beego.Info("Deduct Input: ", deductInput)

	// get service_key from DB with deductInput.service_id
	service, err := models.GetService(deductInput.ServiceID)
	if err != nil {
		b.ResponseError(libs.ErrInvalidService, err)
	}

	var hashed HashedBody
	hashed.ExternalTxID = deductInput.ExternalTxID
	hashed.ExternalItemID = deductInput.ExternalItemID
	hashed.Amount = deductInput.Amount
	hashed.ExternalItemName = deductInput.ExternalItemName
	hashed.ServiceID = deductInput.ServiceID
	bHashed, _ := json.Marshal(hashed)

	// hashed
	h := sha1.New()
	hBody := string(bHashed) + service.Key //
	h.Write([]byte(hBody))
	hashedData := fmt.Sprintf("%x", h.Sum(nil))

	if hashedData != deductInput.Hash {
		beego.Error(hashedData, deductInput.Hash)
		b.ResponseError(libs.ErrInvalidSignature, errors.New(libs.ErrInvalidSignature.Message))
	}

	// test, query with paytransacion and user
	/*
		deductPT, err := models.GetDeductPayTransactionUser(uid)
		if err != nil {
			beego.Error(err)
			b.ResponseError(libs.ErrNoUser, err)
		}
	*/

	// TODO: make log file for inputs... with go routine ??

	// get userinfo for getting balance
	user, err := models.FindByID(uid)
	if err != nil {
		b.ResponseError(libs.ErrNoUser, err)
	}

	// TODO: check external_id in deductHistory. need it ???
	// think more about the transaction !!!
	// reduce three select queries !!!

	// check balance
	//iDeductInputItemAmount, _ := strconv.Atoi(deductInput.ItemAmount)
	//iDeductInputItemAmount := deductInput.Amount

	if user.Balance < deductInput.Amount {
		// low balance
		s := "Need more " + strconv.Itoa(deductInput.Amount-user.Balance) + " balance"
		b.ResponseError(libs.ErrLowBalance, errors.New(s))
	}
	/*
		if deductPT.Balance < deductInput.Amount {
			// low balance
			s := "Need more " + strconv.Itoa(deductInput.Amount-deductPT.Balance) + " balance"
			b.ResponseError(libs.ErrLowBalance, errors.New(s))
		}
	*/

	// check amount_after_used in paytransaction
	//	... amount_after_used 가 0이 아닌것 중에서 가장 오래된 데이터 한건 가져오기...

	deductPT, err := models.GetDeductPayTransaction(user.UID)
	if err != nil {
		b.ResponseError(libs.ErrNoPaytransaction, err)
	}

	// if okay
	//	... deduct, amount_after_used, balance of wallet update.
	// this logic for seperate to paid or free when item sell
	deductFree, deductPaid := 0, 0
	var uf models.UserFilter

	if deductPT.AmountAfterUsed > deductInput.Amount {
		// 해당 paytransaction의 amount_after_used의 amount가 구매 item의 amount 보다 크면, 바로 deduct

		// TODO: logging

		//set kind of deduct
		if deductPT.Price == 0 {
			deductFree = deductInput.Amount
		} else {
			deductPaid = deductInput.Amount
		}

		// calculate
		deductedAmountAfterUsed := libs.Abs(deductPT.AmountAfterUsed - deductInput.Amount)
		//deductedBalance := deductPT.Balance - deductInput.Amount
		deductedBalance := user.Balance - deductInput.Amount

		// make deduct
		//	update amount_after_used in paytransaction
		//	update user_wallet
		uf, err = models.MakeDeduct(deductPT.UID, deductPT.PxID, deductedAmountAfterUsed, deductedBalance)
		if err != nil {
			// TODO: beego error
			b.ResponseError(libs.ErrDatabase, err)
		}
	} else {
		// 해당 paytrasaction의 amount_after_used의 amount가 구매 item의 amount 작으면 ...
		// next paytransacion의 존재 하는 것이다. next paytransaction을 가져오면서 looping 처리 해야 한다.
		nextIDeductInputItemAmount := deductInput.Amount

		// looping
		for nextIDeductInputItemAmount != 0 {
			if nextIDeductInputItemAmount < deductPT.AmountAfterUsed {
				deductHistory := nextIDeductInputItemAmount
				if deductPT.Price == 0 {
					deductFree = deductFree + deductHistory
				} else {
					deductPaid = deductPaid + deductHistory
				}

				// calculate
				deductedAmountAfterUsed := libs.Abs(deductPT.AmountAfterUsed - nextIDeductInputItemAmount)
				//deductedBalance := deductPT.Balance - nextIDeductInputItemAmount
				deductedBalance := user.Balance - nextIDeductInputItemAmount

				uf, err = models.MakeDeduct(deductPT.UID, deductPT.PxID, deductedAmountAfterUsed, deductedBalance)
				if err != nil {
					// TODO: beego error
					b.ResponseError(libs.ErrDatabase, err)
				}

				nextIDeductInputItemAmount = 0

			} else {
				deductHistory := deductPT.AmountAfterUsed
				if deductPT.Price == 0 {
					deductFree = deductFree + deductHistory
				} else {
					deductPaid = deductPaid + deductHistory
				}

				// !!! important, save next amount before update amount_after_used
				nextIDeductInputItemAmount = libs.Abs(nextIDeductInputItemAmount - deductPT.AmountAfterUsed)

				// calculate
				deductedAmountAfterUsed := 0
				//deductedBalance := deductPT.Balance - deductHistory
				deductedBalance := user.Balance - deductHistory

				uf, err = models.MakeDeduct(deductPT.UID, deductPT.PxID, deductedAmountAfterUsed, deductedBalance)
				if err != nil {
					// TODO: beego error
					b.ResponseError(libs.ErrDatabase, err)
				}

				// get next paytransaction for loop
				deductPT, err = models.GetDeductPayTransaction(user.UID)
				//deductPT, err = models.GetDeductPayTransactionUser(uid)
				if err != nil {
					// TODO: beego error
					b.ResponseError(libs.ErrNoPaytransaction, err)
				}

			}
		}

	}

	/* original
	if deductPaytransaction.AmountAfterUsed > iDeductInputItemAmount {
		// 해당 paytransaction의 amount_after_used의 amount가 구매 item의 amount 보다 크면, 바로 deduct

		// TODO: logging

		//set kind of deduct
		if deductPaytransaction.Price == 0 {
			deductFree = iDeductInputItemAmount
		} else {
			deductPaid = iDeductInputItemAmount
		}

		// calculate
		deductedAmountAfterUsed := libs.Abs(deductPaytransaction.AmountAfterUsed - iDeductInputItemAmount)
		deductedBalance := user.Balance - iDeductInputItemAmount

		// make deduct
		//	update amount_after_used in paytransaction
		//	update user_wallet
		uf, err = models.MakeDeduct(user.UID, deductPaytransaction.PxID, deductedAmountAfterUsed, deductedBalance)
		if err != nil {
			// TODO: beego error
			b.ResponseError(libs.ErrDatabase, err)
		}
	} else {
		// 해당 paytrasaction의 amount_after_used의 amount가 구매 item의 amount 작으면 ...
		// next paytransacion의 존재 하는 것이다. next paytransaction을 가져오면서 looping 처리 해야 한다.
		nextIDeductInputItemAmount := iDeductInputItemAmount

		// looping
		for nextIDeductInputItemAmount != 0 {
			if nextIDeductInputItemAmount < deductPaytransaction.AmountAfterUsed {
				deductHistory := nextIDeductInputItemAmount
				if deductPaytransaction.Price == 0 {
					deductFree = deductFree + deductHistory
				} else {
					deductPaid = deductPaid + deductHistory
				}

				// calculate
				deductedAmountAfterUsed := libs.Abs(deductPaytransaction.AmountAfterUsed - nextIDeductInputItemAmount)
				deductedBalance := user.Balance - nextIDeductInputItemAmount

				uf, err = models.MakeDeduct(user.UID, deductPaytransaction.PxID, deductedAmountAfterUsed, deductedBalance)
				if err != nil {
					// TODO: beego error
					b.ResponseError(libs.ErrDatabase, err)
				}

				nextIDeductInputItemAmount = 0

			} else {
				deductHistory := deductPaytransaction.AmountAfterUsed
				if deductPaytransaction.Price == 0 {
					deductFree = deductFree + deductHistory
				} else {
					deductPaid = deductPaid + deductHistory
				}

				// !!! important, save next amount before update amount_after_used
				nextIDeductInputItemAmount = libs.Abs(nextIDeductInputItemAmount - deductPaytransaction.AmountAfterUsed)

				// calculate
				deductedAmountAfterUsed := 0
				deductedBalance := user.Balance - deductHistory

				uf, err = models.MakeDeduct(user.UID, deductPaytransaction.PxID, deductedAmountAfterUsed, deductedBalance)
				if err != nil {
					// TODO: beego error
					b.ResponseError(libs.ErrDatabase, err)
				}

				// get next paytransaction for loop
				deductPaytransaction, err = models.GetDeductPayTransaction(user.UID)
				if err != nil {
					// TODO: beego error
					b.ResponseError(libs.ErrNoPaytransaction, err)
				}

			}
		}

	}
	*/

	//fmt.Println(deductFree, deductPaid)

	// TODO: go routine ???
	// insert deduct history with go routine
	var d models.DeductHistory
	//d.UID = user.UID
	d.UID = deductPT.UID
	d.SID = deductInput.ServiceID
	d.ExternalTxID = deductInput.ExternalTxID
	d.ExternalItemID = deductInput.ExternalItemID
	d.ExternalItemName = deductInput.ExternalItemName
	d.Amount = deductInput.Amount
	d.DeductByFree = deductFree
	d.DeductByPaid = deductPaid
	d.BalanceAfterBuy = uf.Balance

	go models.AddDeductHistory(d)
	/*
		err = models.AddDeductHistory(d)
		if err != nil {
			// TODO: just make a logging file
			beego.Error("AddDeductHistory: ", err)
		}
	*/

	// return
	/*
		var r ResultDeduct
		//d.UID = user.UID
		r.UID = deductPT.UID
		r.ServiceID = deductInput.ServiceID
		r.ItemName = deductInput.ItemName
		r.ItemID = deductInput.ItemID
		r.ItemAmount = deductInput.ItemAmount
		r.ExternalID = deductInput.ExternalID
		r.BalanceAfterBuy = uf.Balance
	*/

	beego.Info("Deducted Time: ", time.Since(start))

	// success
	b.ResponseSuccess("", d)

}

// GetBalance ...
func (b *BillingController) GetBalance() {
	// ...
	// check header authorization.
	// keep time
}

// GetDeductHash ...
// for test or something
func (b *BillingController) GetDeductHash() {
	var input HashedBody

	body, _ := ioutil.ReadAll(b.Ctx.Request.Body)
	err := json.Unmarshal(body, &input)
	if err != nil {
		b.ResponseError(libs.ErrJSONUnmarshal, err)
	}
	// TODO: get auth jwt ???
	// fmt.Println(input)

	// get service_key by service_id
	service, err := models.GetService(input.ServiceID)
	if err != nil {
		b.ResponseError(libs.ErrInvalidService, err)
	}

	bHashed, _ := json.Marshal(input)
	// hashed
	h := sha1.New()
	hBody := string(bHashed) + service.Key //
	h.Write([]byte(hBody))
	hashedData := fmt.Sprintf("%x", h.Sum(nil))

	b.ResponseSuccess("", hashedData)

}

// CallbackXsolla ...
// @Title Get xsolla callback data
// @Description xsolla send callbac data
//
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
		c.TransactionAt = xsollaData.Transaction.PaymentDate

		beego.Info("charge data: ", c)

		// TODO: logging file.

		// begin tran
		err := models.AddPaymentTransaction(c)
		if err != nil {
			beego.Error("Charge error: ", err)
			b.XsollaResponseError(libs.ErrXMakePaytransaction)
		}

		// TODO: set redis?

		// TODO: xsolla success ?
		// success
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
