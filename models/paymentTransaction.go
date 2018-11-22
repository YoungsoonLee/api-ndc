package models

import (
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/orm"
)

// PaymentTransaction ...
type PaymentTransaction struct {
	PxID            string    `orm:"column(PxID);size(500);pk" json:"pxid"`             // paymentTry의 pid
	TxID            string    `orm:"column(TxID);" json:"txid"`                         // pg사로 부터 넘어오는 unique id로 pg사 이용해서 추적이 가능해야 한다.
	UID             string    `orm:"column(UID);size(50);" json:"uid"`                  // user id
	ItemID          int       `orm:"column(ItemID);" json:"itemid"`                     // itemid
	ItemName        string    `orm:"size(1000);" json:"item_name"`                      // not null,
	PgID            int       `orm:"column(PgID);" json:"pgid"`                         // pgid
	Currency        string    `orm:"size(3);default(USD)" json:"currency"`              // not null, default 'USD'
	Price           int       `json:"price"`                                            // not null,
	Amount          int       `json:"amount"`                                           // not null, 실제 적립되는 cyber coin 양
	TransactionAt   time.Time `orm:"type(datetime);auto_now_add" json:"transaction_at"` // 결제 완료일
	AmountAfterUsed int       `json:"amount_after_used"`                                // 사용 후 남은 amount (insert때는 충전되는 amount 와 동일하게...deduct 뙬때 마이너스)
	IsCanceled      bool      `orm:"default(false);null" json:"is_canceled"`            // default: 0(false). 향후 cancel 발생을 대비. 향후 cancel 이력 관련 테이블 필요
	CanceledAt      time.Time `orm:"type(datetime);null" json:"canceled_at"`            // 결제 완료일
}

// DeductPaymentTransaction ...
type DeductPaymentTransaction struct {
	//PaymentTransaction PaymentTransaction
	PxID            string    `orm:"column(PxID);size(500);pk" json:"pxid"`             // paymentTry의 pid
	TxID            string    `orm:"column(TxID);" json:"txid"`                         // pg사로 부터 넘어오는 unique id로 pg사 이용해서 추적이 가능해야 한다.
	UID             string    `orm:"column(UID);size(50);" json:"uid"`                  // user id
	ItemID          int       `orm:"column(ItemID);" json:"itemid"`                     // itemid
	ItemName        string    `orm:"size(1000);" json:"item_name"`                      // not null,
	PgID            int       `orm:"column(PgID);" json:"pgid"`                         // pgid
	Currency        string    `orm:"size(3);default(USD)" json:"currency"`              // not null, default 'USD'
	Price           int       `json:"price"`                                            // not null,
	Amount          int       `json:"amount"`                                           // not null, 실제 적립되는 cyber coin 양
	TransactionAt   time.Time `orm:"type(datetime);auto_now_add" json:"transaction_at"` // 결제 완료일
	AmountAfterUsed int       `json:"amount_after_used"`                                // 사용 후 남은 amount (insert때는 충전되는 amount 와 동일하게...deduct 뙬때 마이너스)
	IsCanceled      bool      `orm:"default(false);null" json:"is_canceled"`            // default: 0(false). 향후 cancel 발생을 대비. 향후 cancel 이력 관련 테이블 필요
	CanceledAt      time.Time `orm:"type(datetime);null" json:"canceled_at"`            // 결제 완료일
	Balance         int       `json:"balance"`
}

// GetPayTransaction ...
func GetPayTransaction(UID string) ([]PaymentTransaction, error) {
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

// AddPaymentTransaction ...
func AddPaymentTransaction(c PaymentTransaction) error {
	o := orm.NewOrm()
	err := o.Begin()

	//_, err = o.Insert(&c)
	sql := "INSERT INTO Payment_transaction " +
		"(\"PxID\", \"TxID\", \"UID\", \"ItemID\", Item_name, \"PgID\", Currency, Price, Amount, Transaction_at, Amount_After_Used) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	_, err = o.Raw(sql, c.PxID, c.TxID, c.UID, c.ItemID, c.ItemName, c.PgID, c.Currency, c.Price, c.Amount, c.TransactionAt, c.Amount).Exec()
	if err != nil {
		beego.Error("AddPaymentTransaction: ", err)
		_ = o.Rollback()
		return err
	}

	sql = "UPDATE \"wallet\" SET balance = balance + ? WHERE \"UID\" = ?"
	_, err = o.Raw(sql, c.Amount, c.UID).Exec()

	if err != nil {
		beego.Error("AddPaymentTransaction update wallet: ", err)
		_ = o.Rollback()
		return err
	}

	err = o.Commit()
	return nil
}

// MakeDeduct ...
func MakeDeduct(UID string, PxID string, deductedAmountAfterUsed, deductedBalance int) (UserFilter, error) {
	o := orm.NewOrm()
	err := o.Begin()

	sql := "UPDATE \"wallet\" SET balance = ? WHERE \"UID\" = ?"
	_, err = o.Raw(sql, deductedBalance, UID).Exec()
	if err != nil {
		beego.Error("Make Deduct Update wallet: ", err)
		_ = o.Rollback()
		return UserFilter{}, err
	}

	// update amount_after_used
	sql = "UPDATE Payment_transaction SET amount_after_used = ? WHERE \"PxID\" = ?"
	_, err = o.Raw(sql, deductedAmountAfterUsed, PxID).Exec()
	if err != nil {
		beego.Error("Make Deduct Update PaymentTransaction: ", err)
		_ = o.Rollback()
		return UserFilter{}, err
	}

	err = o.Commit()

	var user UserFilter
	sql = "SELECT " +
		" \"user\".\"UID\" , " +
		" Displayname, " +
		" Email, " +
		" Confirmed, " +
		" Picture, " +
		" Provider, " +
		" Permission, " +
		" Status, " +
		" \"user\".Create_At, " +
		" \"user\".Update_At, " +
		" \"wallet\".Balance " +
		" FROM \"user\", \"wallet\" " +
		" WHERE \"user\".\"UID\" = \"wallet\".\"UID\" " +
		" and \"user\".\"UID\" = ?"
	err = o.Raw(sql, UID).QueryRow(&user)
	return user, nil
}

// GetDeductPayTransaction ...
func GetDeductPayTransaction(UID string) (PaymentTransaction, error) {
	o := orm.NewOrm()
	var pt PaymentTransaction

	// ... amount_after_used 가 0이 아닌것 중에서 가장 오래된 데이터 한건 가져오기...
	sql := " SELECT \"PxID\", \"TxID\", \"UID\", \"ItemID\", Item_name, \"PgID\", " +
		" Currency, Price, Amount, Transaction_at, Amount_After_Used " +
		" FROM Payment_transaction " +
		" WHERE \"UID\" = ? " +
		" and Amount_After_Used != 0 " +
		" ORDER by Transaction_at asc LIMIT 1"
	err := o.Raw(sql, UID).QueryRow(&pt)
	if err != nil {
		return PaymentTransaction{}, err
	}

	return pt, nil
}

// GetDeductPayTransactionUser ...
func GetDeductPayTransactionUser(UID string) (DeductPaymentTransaction, error) {
	o := orm.NewOrm()
	var ptu DeductPaymentTransaction

	// ... amount_after_used 가 0이 아닌것 중에서 가장 오래된 데이터 한건 가져오기...
	sql := " SELECT \"PxID\", \"TxID\", Payment_transaction.\"UID\", \"ItemID\", Item_name, \"PgID\", " +
		" Currency, Price, Amount, Transaction_at, Amount_After_Used, " +
		" \"wallet\".balance as balance" +
		" FROM Payment_transaction , \"wallet\" " +
		" WHERE Payment_transaction.\"UID\" = \"wallet\".\"UID\" " +
		" and \"wallet\".\"UID\" = ? " +
		" and Amount_After_Used != 0 " +
		" ORDER by Transaction_at asc LIMIT 1"
	err := o.Raw(sql, UID).QueryRow(&ptu)
	if err != nil {
		return DeductPaymentTransaction{}, err
	}

	return ptu, nil
}
