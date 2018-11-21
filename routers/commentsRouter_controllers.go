package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"],
		beego.ControllerComments{
			Method: "CheckDisplayName",
			Router: `/:displayname`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"],
		beego.ControllerComments{
			Method: "CreateUser",
			Router: `/CreateUser`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Social",
			Router: `/Social`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:BillingController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:BillingController"],
		beego.ControllerComments{
			Method: "GetChargeItems",
			Router: `/`,
			AllowHTTPMethods: []string{"GET"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:BillingController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:BillingController"],
		beego.ControllerComments{
			Method: "GetPaymentToken",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:PaymentCategoryController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:PaymentCategoryController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:PaymentGatewayController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:PaymentGatewayController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:PaymentItemController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:PaymentItemController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:ServiceController"] = append(beego.GlobalControllerRouter["github.com/YoungsoonLee/api-ndc/controllers:ServiceController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
