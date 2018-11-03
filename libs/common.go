package libs

// ControllerError is controller error info structer.
type ControllerError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	//DevInfo  string `json:"dev_info"`
	//MoreInfo string `json:"more_info"`
}

// Predefined controller error values.
var (
	Err404 = &ControllerError{404, "404", "page not found"}

	// 10000 ~ related on account or auth
	ErrInputData        = &ControllerError{400, "10001", "Data input error"}
	ErrDisplayname      = &ControllerError{400, "10002", "Displayname should have 4 ~ 16 letters."}
	ErrEmail            = &ControllerError{400, "10003", "Must be a valid email address"}
	ErrMaxEmail         = &ControllerError{400, "10004", "Displayname cannot have over 100 letters."}
	ErrPassword         = &ControllerError{400, "10005", "Password should have 8 ~ 16 letters with number and special characters"}
	ErrDupDisplayname   = &ControllerError{400, "10006", "Displayname already exists"}
	ErrDupEmail         = &ControllerError{400, "10007", "Email already exists"}
	ErrAlreadyConfirmed = &ControllerError{400, "10008", "Email already confirmed."}
	ErrWrongToken       = &ControllerError{400, "10009", "wrong token."}
	ErrExpiredToken     = &ControllerError{401, "10010", "The token was already expired or invalid token. try again."}
	ErrPass             = &ControllerError{400, "10011", "User information does not exist or the password is incorrect"}
	ErrTokenAbsent      = &ControllerError{400, "10012", "Token absent"}
	ErrTokenInvalid     = &ControllerError{400, "10013", "Token invalid"}
	ErrTokenOther       = &ControllerError{400, "10014", "Token other"}
	ErrNoUser           = &ControllerError{400, "10015", "User information does not exist"}
	ErrIDAbsent         = &ControllerError{400, "10016", "Id absent"}
	ErrLoginFacebook    = &ControllerError{400, "10017", "Your disaplayname is connected a facebook. use facebook login."}
	ErrLoginGoogle      = &ControllerError{400, "10018", "Your disaplayname is connected a Google. use Google login."}

	/*
		ErrNoUserPass   = &ControllerError{400, "10006", "User information does not exist or the password is incorrect"}
		ErrNoUserChange = &ControllerError{400, "10007", "User information does not exist or data has not changed"}
		ErrInvalidUser  = &ControllerError{400, "10008", "User information is incorrect"}
		ErrOpenFile     = &ControllerError{500, "10009", "Error opening file"}
		ErrWriteFile    = &ControllerError{500, "10010", "Error writing a file"}
		ErrSystem       = &ControllerError{500, "10011", "Operating system error"}
		ErrExpired      = &ControllerError{400, "10012", "Login has expired"}
		ErrPermission   = &ControllerError{400, "10013", "Permission denied"}
	*/

	// 20000 ~ related in payment
	ErrNoPaymentItem    = &ControllerError{400, "20001", "PaymentItem information does not exists"}
	ErrNoCategoryID     = &ControllerError{400, "20002", "PaymentCategory information does not exists"}
	ErrNoPGID           = &ControllerError{400, "20003", "PaymentGateway information does not exists"}
	ErrNoSignature      = &ControllerError{400, "20004", "Deduct Signature does not exists"}
	ErrInvalidSignature = &ControllerError{400, "20005", "Invalid Signature"}
	ErrInvalidService   = &ControllerError{400, "20006", "Invalid Service"}

	// 90000 ~ related on system error
	ErrDatabase      = &ControllerError{500, "90001", "Database operation error"}
	ErrJSONUnmarshal = &ControllerError{500, "90002", "JSON Unmarshal error"}
	ErrFailSalt      = &ControllerError{500, "90003", "Generate error"}
	ErrFailHash      = &ControllerError{500, "90004", "Generate Hash error"}
	ErrMakeToken     = &ControllerError{500, "90005", "Generate Token error"}
	ErrJSONmarshal   = &ControllerError{500, "90006", "JSON marshal error"}
	ErrTokenRequest  = &ControllerError{500, "90007", "Token Request error with xsolla"}
	ErrClient        = &ControllerError{500, "90008", "Get client data error with xsolla"}

	// xsolla
	ErrXNilSig             = &ControllerError{400, "INVALID_SIGNATURE_SIGNATURE_NULL", "INVALID_SIGNATURE_SIGNATURE_NULL"}
	ErrXInvalidSig         = &ControllerError{400, "INVALID_SIGNATURE", "INVALID_SIGNATURE"}
	ErrXInvalidUser        = &ControllerError{400, "INVALID_USER", "INVALID_USER"}
	ErrXInvalidJSON        = &ControllerError{400, "JSON_PARSING_ERROR", "JSON_PARSING_ERROR"}
	ErrXInvalidPaytryData  = &ControllerError{400, "INVALID_PAYTRY_DATA", "INVALID_PAYTRY_DATA"}
	ErrXMakePaytransaction = &ControllerError{400, "ERROR_MAKE_PAYTRANSACTION", "ERROR_MAKE_PAYTRANSACTION"}
	ErrXInvalidNotiType    = &ControllerError{400, "INVALID_NOTI_TYPE", "INVALID_NOTI_TYPE"}
)
