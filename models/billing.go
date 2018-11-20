package models

import (
	"github.com/astaxie/beego/orm"
)

// GetChargeItems ...
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

// GetPayTransaction ...
func GetPayTransaction(UID int64) ([]PaymentTransaction, error) {
	var payTransactions []PaymentTransaction

	o := orm.NewOrm()
	sql := "SELECT " +
		" \"PxID\" , " +
		" \"TxID\", " +
		" Item_Name, " +
		" Price, " +
		" Amount, " +
		" Transaction_At" +
		" FROM \"payment_transaction\" " +
		" WHERE \"UID\" = ? " +
		" ORDER BY Transaction_At desc "
	_, err := o.Raw(sql, UID).QueryRows(&payTransactions)
	return payTransactions, err
}

// GetUsedHistory ...
func GetUsedHistory(UID int64) ([]DeductHistory, error) {
	var deductHistory []DeductHistory

	o := orm.NewOrm()
	sql := "SELECT " +
		" \"ID\" , " +
		" \"UID\", " +
		" \"ExternalID\", " +
		" Item_Name, " +
		" Amount, " +
		" Deduct_by_free, " +
		" Deduct_by_paid, " +
		" Used_at " +
		" FROM \"deduct_history\" " +
		" WHERE \"UID\" = ? " +
		" ORDER BY Used_at desc "
	_, err := o.Raw(sql, UID).QueryRows(&deductHistory)
	return deductHistory, err
}
