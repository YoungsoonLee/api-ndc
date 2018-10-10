package controllers

import (
	"encoding/json"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
)

type PaymentItemController struct {
	BaseController
}

// Post ...
// @Title Create Payment Category
// @Description create payment category
// @Param	categoryid			json 	INT		false		"payment category id"
// @Param	item_name			json 	string	false		"item name"
// @Param	item_description	json 	string	false		"item description"
// @Param	pgid				json 	string	false		"payment gateway id"
// @Param	currency			json 	string	false		"currency"
// @Param	price				json 	INT		false		"price"
// @Param	amount				json 	INT		false		"amount of charge cyber coin"
// @Success 200 {int} models.PaymentItem.ItemID
// @Failure 403 body is empty
// @router / [post]
func (p *PaymentItemController) Post() {

	var pi models.PaymentItem

	err := json.Unmarshal(p.Ctx.Input.RequestBody, &pi)
	if err != nil {
		p.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// TODO: validation

	// save to db
	itemid, err := models.AddPaymentItem(pi)
	if err != nil {
		p.ResponseError(libs.ErrDatabase, err)
	}

	//success
	p.ResponseSuccess("ItemID", itemid)
}
