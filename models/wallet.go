package models

import "time"

// Wallet ...
type Wallet struct {
	UID      int64     `orm:"column(UID);pk"`              // user id
	Balance  int       `orm:"default(0)"`                  // balances
	CreateAt time.Time `orm:"type(datetime);auto_now_add"` // first save
	UpdateAt time.Time `orm:"type(datetime);auto_now"`     // eveytime save
}
