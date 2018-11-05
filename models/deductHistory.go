package models

import "time"

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
	ID           int64     `orm:"auto;pk" json:"id"`                                // DeductHistory id
	UID          int64     `orm:"column(UID);" json:"uid"`                          // user id
	SID          string    `orm:"column(SID);size(500);" json:"sid"`                //
	ExternalID   string    `orm:"column(ExternalID);size(500);" json:"external_id"` // 아이템 구매시 각 게임 서버로 부터 오는 고유의 트랜잭션 ID
	ItemID       int       `orm:"column(ItemID);" json:"itemid"`                    // 각 게임별 고유의 item_id 혹은 고유의 추적이 가능한 무엇
	ItemName     string    `orm:"size(1000);" json:"item_name"`                     // item 이름
	Amount       int       `json:"amount"`                                          // 아이템의 deduct coin 양
	DeductByFree int       `json:"deduct_by_free"`                                  // 무료 사이버머니로 구입된 양
	DeductByPaid int       `json:"deduct_by_paid"`                                  // 유료 사이버머니로 구입된 양
	UsedAt       time.Time `orm:"type(datetime);" json:"used_at"`                   // 사용 일
	IsCanceled   bool      `orm:"default(false);null" json:"is_canceled"`           // default: 0(false). 향후 cancel 발생을 대비. 향후 cancel 이력 관련 테이블 필요
	CanceledAt   time.Time `orm:"type(datetime);null" json:"canceled_at"`           //
}
