package libs

type ID struct {
	Value  int64 `json:"value"`
	Hidden bool  `json:"hidden"`
}

type Email struct {
	Value       string `json:"value"`
	AllowModify bool   `json:"allow_modify"`
	Hidden      bool   `json:"hidden"`
}

type Country struct {
	Value string `json:"value"`
}

type Name struct {
	Value  string `json:"value"`
	Hidden bool   `json:"hidden"`
}

type User struct {
	ID      ID      `json:"id"`
	Email   Email   `json:"email"`
	Country Country `json:"country"`
	Name    Name    `json:"name"`
}

type UI struct {
	Size string `json:"size"`
}

type Settings struct {
	ProjectID  int    `json:"project_id"`
	ExternalID string `json:"external_id"`
	Mode       string `json:"mode"`
	Language   string `json:"language"`
	Currency   string `json:"currency"`
	UI         UI     `json:"ui"`
}

type Checkout struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}

type Description struct {
	value string `json:"value"`
}

type Purchase struct {
	Checkout    Checkout    `json:"checkout"`
	Description Description `json:"description"`
}

type CustomParameters struct {
	Pid string `json:"pid"`
}

type XsollaSendJSONToGetToken struct {
	User             User             `json:"user"`
	Settings         Settings         `json:"settings"`
	Purchase         Purchase         `json:"purchase"`
	CustomParameters CustomParameters `json:"custom_parameters"`
}
