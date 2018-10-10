package controllers

import (
	"encoding/json"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
)

type PaymentCategoryController struct {
	BaseController
}

// Post ...
// @Title Create Payment Category
// @Description create payment category
// @Param	category				json 	INT		false		"100: for paid charge, 200: free rewards, 300: free bonus"
// @Param	category_description	json 	string	false		"category description"
// @Success 200 {int} models.PaymentCategory.CategoryID
// @Failure 403 body is empty
// @router / [post]
func (p *PaymentCategoryController) Post() {

	var pc models.PaymentCategory

	err := json.Unmarshal(p.Ctx.Input.RequestBody, &pc)
	if err != nil {
		p.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// TODO: validation

	// save to db
	pcid, err := models.AddPaymentCategory(pc)
	if err != nil {
		p.ResponseError(libs.ErrDatabase, err)
	}

	//success
	p.ResponseSuccess("pcid", pcid)
}
