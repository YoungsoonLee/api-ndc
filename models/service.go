package models

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Service struct {
	SID         string    `orm:"column(SID);size(500);pk" json:"sid"`          // service id
	Key         string    `orm:"size(500);unique" json:"key"`                  // key for encrypt
	Description string    `orm:"size(500)" json:"description"`                 //
	CreateAt    time.Time `orm:"type(datetime);auto_now_add" json:"create_at"` // first save
	CloseAt     time.Time `orm:"type(datetime);null" json:"CloseAt"`           // eveytime save
}

func AddService(s Service) (string, error) {

	b := make([]byte, 8) //equals 16 charachters
	rand.Read(b)

	// make Id
	s.SID = "S" + strconv.FormatInt(time.Now().UnixNano(), 10)
	s.Key = hex.EncodeToString(b)

	_, err := orm.NewOrm().Raw("INSERT INTO service (\"SID\", key, Description, Create_At) VALUES ($1, $2, $3, $4)", s.SID, s.Key, s.Description, time.Now()).Exec()
	if err != nil {
		return "", err
	}

	return s.SID, nil
}
