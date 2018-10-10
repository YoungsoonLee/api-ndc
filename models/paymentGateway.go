package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
  payment_gateway              //PG사 정보 관리 테이블
      pg_id                       // unique, not auto increse
      pg_description              //
      pg_kind                     // 향후 사용 할 수 있다. ex) 1: credit card. 2: mobile ....
      created_at                  //
      closed_at                   //
*/
type PaymentGateway struct {
	PgID          int       `orm:"column(PgID);pk;auto" json:"pgid"`             // unique, auto increase
	PgDescription string    `orm:"size(1000);" json:"pg_description"`            // not null,
	PgKind        int       `orm:"null;" json:"pg_kind"`                         // 향후 사용 할 수 있다. ex) 1: credit card. 2: mobile ....
	CreateAt      time.Time `orm:"type(datetime);auto_now_add" json:"create_at"` //
	CloseAt       time.Time `orm:"type(datetime);null" json:"close_at"`          //
}

func AddPaymentGateway(pg PaymentGateway) (int, error) {

	_, err := orm.NewOrm().Insert(&pg)
	if err != nil {
		beego.Error("Error AddPaymentGateway: ", err)
		return 0, err
	}

	return pg.PgID, nil
}
