package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

/**
  payment_transaction          // 결제 완료 테이블(charge). 파라미터 비교시 payment_try 내용과 다를경우 hacking으로 간주
     PxID                         // payment_try의 pid
     transaction_id              // pg사로 부터 넘어오는 unique id로 pg사 이용해서 추적이 가능해야 한다.
     user_id
     item_id
     pg_id
     currency                    // default: 'USD'.
     price
     amount                      //cyber coin amount
     transaction_at              // 결제 완료일
     amount_after_used           // 사용 후 남은 amount (insert시 충전되는 amount 와 동일하게...deduct 뙬때 마이너스)
     is_canceled                 // default: 0(false). 향후 cancel 발생을 대비. 향후 cancel 이력 관련 테이블 필요
     canceled_at                //
*/
type PaymentTransaction struct {
	PxID            string    `orm:"column(PxID);size(500);pk" json:"pxid"`             // paymentTry의 pid
	TxID            string    `orm:"column(TxID);" json:"txid"`                         // pg사로 부터 넘어오는 unique id로 pg사 이용해서 추적이 가능해야 한다.
	UID             int64     `orm:"column(UID);" json:"uid"`                           // user id
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

func AddPaymentTransaction(c PaymentTransaction) error {
	o := orm.NewOrm()
	err := o.Begin()

	_, err = o.Insert(&c)
	if err != nil {
		err = o.Rollback()
		return err
	}

	var wallet Wallet
	sql := "UPDATE wallet set balance = sum(balance) + ? WHERE \"UID\" = ?"
	err = o.Raw(sql, c.Amount, c.UID).QueryRow(&wallet)
	if err != nil {
		err = o.Rollback()
	}

	err = o.Commit()
	return nil

}
