package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
  user_deduct_history     // coin의 차감이 발생 할때 기록 하는 테이블. ex)게임 아이템 구매시
     deduct_id                // unique,  auto increase
     user_id
     service_id
     external_id              // 아이템 구매시 각 게임 서버로 부터 오는 고유의 트랜잭션 ID
     item_id                  // 각 게임별 고유의 item_id 혹은 고유의 추적이 가능한 무엇
     item_name                // item 이름
     item_amount              // 아이템의 deduct coin 양
     deduct_by_free           // 무료 사이버머니로 구입된 양
     deduct_by_paid           // 유료 사이버머니로 구입된 양
     used_at                  //
     is_canceled              // default: 0(false). 향후 cancel 발생을 대비. 향후 cancel 이력 관련 테이블 필요
     canceled_at              //
*/
type DeductHistory struct {
	ID               int64     `orm:"column(ID);auto;pk" json:"id"`                              // DeductHistory id
	UID              string    `orm:"column(UID);size(50);" json:"uid"`                          // user id
	SID              string    `orm:"column(SID);size(500);" json:"sid"`                         //
	ExternalTxID     string    `orm:"column(ExternalTxID);size(500);" json:"external_txid"`      // 아이템 구매시 각 게임 서버로 부터 오는 고유의 트랜잭션 ID
	ExternalItemID   string    `orm:"column(ExternalItemID);size(1000);" json:"external_itemid"` // 각 게임별 고유의 item_id 혹은 고유의 추적이 가능한 무엇
	ExternalItemName string    `orm:"size(1000);" json:"external_itemname"`                      // 각 게임별 item 이름
	Amount           int       `json:"amount"`                                                   // 아이템의 deduct coin 양
	DeductByFree     int       `json:"deduct_by_free"`                                           // 무료 사이버머니로 구입된 양
	DeductByPaid     int       `json:"deduct_by_paid"`                                           // 유료 사이버머니로 구입된 양
	UsedAt           time.Time `orm:"type(datetime);auto_now_add" json:"used_at"`                // 사용 일
	IsCanceled       bool      `orm:"default(false);null" json:"is_canceled"`                    // default: 0(false). 향후 cancel 발생을 대비. 향후 cancel 이력 관련 테이블 필요
	CanceledAt       time.Time `orm:"type(datetime);null" json:"canceled_at"`                    //
	BalanceAfterBuy  int       `orm:"-" json:"balance_after_buy"`
}

// GetUsedHistory ...
func GetUsedHistory(UID string) ([]DeductHistory, error) {
	var deductHistory []DeductHistory

	o := orm.NewOrm()
	sql := "SELECT " +
		" \"ID\" , " +
		" \"UID\", " +
		" \"ExternalTxID\", " +
		" \"ExternalItemID\", " +
		" external_item_name, " +
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

// AddDeductHistory ...
func AddDeductHistory(d DeductHistory) error {
	o := orm.NewOrm()

	//_, err = o.Insert(&c)
	sql := "INSERT INTO Deduct_history " +
		"(\"UID\", \"SID\", \"ExternalTxID\", \"ExternalItemID\", external_item_name, amount, deduct_by_free, deduct_by_paid, used_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp)"

	_, err := o.Raw(sql, d.UID, d.SID, d.ExternalTxID, d.ExternalItemID, d.ExternalItemName, d.Amount, d.DeductByFree, d.DeductByPaid).Exec()
	if err != nil {
		beego.Error("AddDeductHistory: ", err)
		return err
	}
	return nil
}
