package controllers

import (
	"encoding/json"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"
)

type ServiceController struct {
	BaseController
}

// Post ...
// @Title CreateService
// @Description create services
// @Param	body		body 	models.Service	true		"body for service content"
// @Success 200 {int} models.Service.Id
// @Failure 403 body is empty
// @router / [post]
func (s *ServiceController) Post() {

	var service models.Service
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &service)
	if err != nil {
		s.ResponseError(libs.ErrJSONUnmarshal, err)
	}

	// TODO: validation

	// save to db
	sid, err := models.AddService(service)
	if err != nil {
		s.ResponseError(libs.ErrDatabase, err)
	}

	//success
	s.ResponseSuccess("sid", sid)
}
