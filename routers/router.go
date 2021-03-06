// @APIVersion 1.0.0
// @Title Naddic platform API
// @Description Naddic platform API
// @Contact youngtip@naddic.com
// @TermsOfServiceUrl
// @License
// @LicenseUrl
package routers

import (
	"github.com/YoungsoonLee/api-ndc/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/user",
			beego.NSRouter("/confirmEmail/:confirmToken", &controllers.UserController{}, "post:ConfirmEmail"),
			beego.NSRouter("/resendConfirmEmail/:email", &controllers.UserController{}, "post:ResendConfirmEmail"),
			beego.NSRouter("/forgotPassword/:email", &controllers.UserController{}, "post:ForogtPassword"),
			beego.NSRouter("/isValidResetPasswordToken/:resetToken", &controllers.UserController{}, "post:IsValidResetPasswordToken"),
			beego.NSRouter("/resetPassword/", &controllers.UserController{}, "post:ResetPassword"),

			beego.NSRouter("/getProfile", &controllers.UserController{}, "post:GetProfile"),
			beego.NSRouter("/updateProfile/", &controllers.UserController{}, "post:UpdateProfile"),

			beego.NSRouter("/updatePassword/", &controllers.UserController{}, "post:UpdatePassword"),
		),

		beego.NSNamespace("/auth",
			beego.NSRouter("/checkDisplayName/:displayname", &controllers.AuthController{}, "get:CheckDisplayName"),
			beego.NSRouter("/register", &controllers.AuthController{}, "post:CreateUser"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
			beego.NSRouter("/checkLogin", &controllers.AuthController{}, "get:CheckLogin"),
			beego.NSRouter("/social", &controllers.AuthController{}, "post:Social"),
			//beego.NSRouter("/logout", &controllers.AuthController{}, "post:Logout"),
		),

		beego.NSNamespace("/billing",
			beego.NSRouter("/getChargeItems", &controllers.BillingController{}, "get:GetChargeItems"),
			beego.NSRouter("/getPaymentToken", &controllers.BillingController{}, "post:GetPaymentToken"),
			beego.NSRouter("/callbackXsolla", &controllers.BillingController{}, "post:CallbackXsolla"),
			beego.NSRouter("/getChargeHistory/:UID", &controllers.BillingController{}, "post:GetChargeHistory"),
			beego.NSRouter("/getUsedHistory/:UID", &controllers.BillingController{}, "post:GetUsedHistory"),
			beego.NSRouter("/buyItem", &controllers.BillingController{}, "post:BuyItem"),
			beego.NSRouter("/getDeductHash", &controllers.BillingController{}, "post:GetDeductHash"),
			beego.NSRouter("/getBalance", &controllers.BillingController{}, "post:GetBalance"),
			//beego.NSRouter("/testBuyItem", &controllers.BillingController{}, "post:TestBuyItem"),
		),

		// news

		//adimn
		beego.NSNamespace("/admin",
			beego.NSNamespace("/service", beego.NSInclude(&controllers.ServiceController{})),
			beego.NSNamespace("/paymentgateway", beego.NSInclude(&controllers.PaymentGatewayController{})),
			beego.NSNamespace("/paymentcategory", beego.NSInclude(&controllers.PaymentCategoryController{})),
			beego.NSNamespace("/paymentitem", beego.NSInclude(&controllers.PaymentItemController{})),
		),
	)
	beego.AddNamespace(ns)
}
