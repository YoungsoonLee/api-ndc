package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
  payment_category            //아레 item이 어느 카데고리인지 관리 하는 테이블. 유료 무료, 무료일 경우 어떠한 무료 인지 구분 짓는 테이블 이다.(주로 통계용이 주목적이다.)
     cid                         // unique, auto increse. pk
     category_id                 // 100: 유료 충전용, 200: 무료 rewars, 300: 무료 bonus
     category_description        //
     sub_category_id             // 0:paid(charge paid coin), ex)category_id를 다시 상세화 할 때 사용. 주로 통계용
     sub_category_description    //
     created_at                  // .defaultTo(knex.fn.now())
     closed_at                   //         [description]
*/
type PaymentCategory struct {
	CategoryID             int       `orm:"column(CategoryID);pk;auto" json:"categoryid"`    // unique, auto increase
	Category               int       `json:"category"`                                       // 100: 유료 충전용, 200: 무료 rewards, 300: 무료 bonus
	CategoryDescription    string    `orm:"size(1000);" json:"category_description"`         // not null,
	SubCategory            int       `orm:"null;" json:"sub_category"`                       //
	SubCategoryDescription string    `orm:"size(1000);null" json:"sub_category_description"` // ex)category_id를 다시 상세화 할 때 사용. 주로 통계용
	CreateAt               time.Time `orm:"type(datetime);auto_now_add" json:"create_at"`    //
	CloseAt                time.Time `orm:"type(datetime);null" json:"close_at"`             //
}

func AddPaymentCategory(pc PaymentCategory) (int, error) {
	_, err := orm.NewOrm().Insert(&pc)
	if err != nil {
		beego.Error("Error AddPaymentCategory: ", err)
		return 0, err
	}

	return pc.CategoryID, nil
}
