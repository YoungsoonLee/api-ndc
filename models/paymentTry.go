package models

import (
	"crypto/rand"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/YoungsoonLee/api-ndc/libs"
	"github.com/astaxie/beego/orm"
)

/**
  payment_try                  // 결제 시도 테이블. PG사 이용시에 이용한다. ex) PG사 오픈전에 기록된다.
     pid                         // unique, not auto increase. make uid like 'P******'. string
     user_id
     item_id
     pg_id
     currency                    // default: 'USD'.
     price
     amount
     tried_at                    // 결제 시도 일자
*/
type PaymentTry struct {
	PxID     string    `orm:"column(PxID);size(500);pk" json:"pxid"`       // unique, payment transaction id
	UID      int64     `orm:"column(UID);" json:"uid"`                     // user id
	ItemID   int       `orm:"column(ItemID);" json:"itemid"`               // itemid
	ItemName string    `orm:"size(1000);" json:"item_name"`                // not null,
	PgID     int       `orm:"column(PgID);" json:"pgid"`                   // pgid
	Currency string    `orm:"size(3);default(USD)" json:"currency"`        // not null, default 'USD'
	Price    int       `json:"price"`                                      // not null,
	Amount   int       `json:"amount"`                                     // not null, 실제 적립되는 cyber coin 양
	TriedAt  time.Time `orm:"type(datetime);auto_now_add" json:"tried_at"` // first save
	Mode     string    `orm:"-" json:"mode"`                               // xsolla mode
	Token    string    `orm:"-" json:"token"`                              // xsolla token
}

// AddPaymentTry ...
func AddPaymentTry(pt PaymentTry) (PaymentTry, error) {
	// check UID
	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("UID", pt.UID).Exist()
	if !exist {
		return PaymentTry{}, errors.New(libs.ErrNoUser.Message)
	}

	// set PgID, Currency, Price, Amount through paymentItem
	sql := "SELECT \"ItemID\", \"PgID\", Item_name, Currency, Price, Amount FROM Payment_Item WHERE \"ItemID\" = ?"
	err := o.Raw(sql, pt.ItemID).QueryRow(&pt)
	if err != nil {
		return PaymentTry{}, err
	}

	// set PxID
	b := make([]byte, 8) //equals 16 charachters
	rand.Read(b)
	pt.PxID = "Px" + strconv.FormatInt(time.Now().UnixNano(), 10)

	sql = "INSERT INTO payment_try" +
		" (\"PxID\", \"UID\", \"ItemID\", \"PgID\", Item_name, Currency, Price, Amount, Tried_At)" +
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err = o.Raw(sql, pt.PxID, pt.UID, pt.ItemID, pt.PgID, pt.ItemName, pt.Currency, pt.Price, pt.Amount, time.Now()).Exec()
	if err != nil {
		return PaymentTry{}, err
	}

	/*
		if beego.BConfig.RunMode == "dev" {
			pt.Mode = "sandbox"
		} else {
			pt.Mode = "production"
		}
	*/
	pt.Mode = os.Getenv("XSOLLA_MODE")

	return pt, nil
}

func CheckPaymentTry(pt PaymentTry) (PaymentTry, bool) {
	o := orm.NewOrm()
	err := o.QueryTable("PaymentTry").Filter("PxID", pt.PxID).Filter("UID", pt.UID).Filter("Amount", pt.Amount).One(&pt)
	if err == orm.ErrMultiRows || err == orm.ErrNoRows {
		return PaymentTry{}, false
	}

	return pt, true

}
