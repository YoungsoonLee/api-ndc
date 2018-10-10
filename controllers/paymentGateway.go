package controllers

import (
	"encoding/json"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
)

type PaymentGatewayController struct {
	BaseController
}

// Post ...
// @Title Create Payment Gateway
// @Description create payment gateway
// @Param	pg_description	json 	string	false		"pg description"
// @Success 200 {int} models.PaymentGateway.PgID
// @Failure 403 body is empty
// @router / [post]
func (p *PaymentGatewayController) Post() {

	var pg models.PaymentGateway
	err := json.Unmarshal(p.Ctx.Input.RequestBody, &pg)
	if err != nil {
		p.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// TODO: validation

	// save to db
	pgid, err := models.AddPaymentGateway(pg)
	if err != nil {
		p.ResponseError(libs.ErrDatabase, err)
	}

	//success
	p.ResponseSuccess("pgid", pgid)
}
