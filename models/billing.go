package models

import (
	"github.com/astaxie/beego/orm"
)

// AddPaymentItem ...
// TODO: you need pgid ???
func GetChargeItems() ([]PaymentItem, error) {
	var chargeItems []PaymentItem

	o := orm.NewOrm()
	sql := "SELECT * FROM Payment_Item WHERE Close_at is null" // close is null
	_, err := o.Raw(sql).QueryRows(&chargeItems)
	if err != nil {
		return chargeItems, err
	}

	return chargeItems, nil
}
