package controllers

import (
	"regexp"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/YoungsoonLee/api-ndc/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

// BaseController ...
type BaseController struct {
	beego.Controller
}

// ResponseError ...
//func (b *BaseController) ResponseError(code string, err error) {
func (b *BaseController) ResponseError(e *libs.ControllerError, err error) {
	// TODO: logging
	devInfo := ""
	if err != nil {
		devInfo = err.Error()
	}
	beego.Error(b.Ctx.Request.RequestURI, e.Message, devInfo)

	response := &models.RespCode{
		Code:    e.Code,
		Message: e.Message,
		DevInfo: devInfo,
		Data:    nil,
	}
	b.Ctx.Output.Status = e.Status
	b.Ctx.Output.JSON(response, true, true)

	// TODO: logging
	b.StopRun()
}

func (b *BaseController) XsollaResponseError(e *libs.ControllerError) {
	// TODO: logging
	beego.Error(b.Ctx.Request.RequestURI, e.Message)

	eData := models.XRespDetailCode{
		Code:    e.Code,
		Message: e.Message,
	}

	response := &models.XRespCode{
		Error: eData,
	}

	b.Ctx.Output.Status = e.Status
	b.Ctx.Output.JSON(response, true, true)

	// TODO: logging
	b.StopRun()
}

// ValidDisplayname ...
func (b *BaseController) ValidDisplayname(displayname string) {

	if len(displayname) < 4 || len(displayname) > 16 {
		//beego.Error("key: displayname, value: ", displayname, ", message: ", libs.ErrDisplayname.Message)
		b.ResponseError(libs.ErrDisplayname, nil)
	}
}

// ValidID ...
func (b *BaseController) ValidID(id string) {

	if len(id) == 0 {
		b.ResponseError(libs.ErrIDAbsent, nil)
	}
}

// ValidEmail ...
func (b *BaseController) ValidEmail(email string) {
	valid := validation.Validation{}
	v := valid.Email(email, "Email")
	if !v.Ok {
		//loggingValidError(v)
		b.ResponseError(libs.ErrEmail, nil)
	}

	v = valid.MaxSize(email, 100, "Email")
	if !v.Ok {
		//loggingValidError(v)
		b.ResponseError(libs.ErrMaxEmail, nil)
	}
}

// ValidPassword ...
func (b *BaseController) ValidPassword(password string) {
	// 8 ~ 16 letters
	if len(password) < 8 || len(password) > 16 {
		b.ResponseError(libs.ErrPassword, nil)
	}

	valid := validation.Validation{}
	pattern := regexp.MustCompile("") //TODO: add regex for password

	v := valid.Match(password, pattern, "password")
	if !v.Ok {
		b.ResponseError(libs.ErrPassword, nil)
	}
}

/*
// ResponseHTTPError ...
func (b *BaseController) ResponseHTTPError(status int, code string, err error) {
	b.Ctx.Output.Status = status
	b.ResponseError(code, err)
}

// ResponseCommonError ...
func (b *BaseController) ResponseCommonError(e *libs.ControllerError) {
	beego.Error(fmt.Errorf(e.Message))
	b.ResponseHTTPError(e.Status, e.Code, fmt.Errorf(e.Message))
}

// ResponseServerError ...
func (b *BaseController) ResponseServerError(e *libs.ControllerError, err error) {
	beego.Error(err)
	b.ResponseHTTPError(e.Status, e.Code, fmt.Errorf(e.Message))
}

func loggingValidError(v *validation.Result) {
	beego.Error("key: ", v.Error.Key, ", value: ", v.Error.Value, ", message: ", v.Error.Message)
}
*/

// ResponseSuccess ...
func (b *BaseController) ResponseSuccess(key string, value interface{}) {
	b.Ctx.Output.Status = 200

	if key == "" {
		mresponse := &models.MrespCode{
			Code:    "ok",
			Message: "success",
			Data:    value,
		}
		b.Ctx.Output.JSON(mresponse, true, true)
		//b.StopRun()
	}

	if key == "tabulator" {
		/*
			tresponse := &models.MrespCode{
				Data: value,
			}
		*/

		b.Ctx.Output.JSON(value, true, true)
		//b.StopRun()
	}

	response := &models.RespCode{
		Code:    "ok",
		Message: "success",
		Data:    map[string]interface{}{},
	}

	response.Data[key] = value
	b.Ctx.Output.JSON(response, true, true)
	//b.StopRun()

}
