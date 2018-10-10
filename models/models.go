package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	//_ "github.com/go-sql-driver/mysql"
)

func RegisterDB() {
	// register model
	orm.RegisterModel(new(User), new(Service), new(Wallet), new(PaymentGateway), new(PaymentCategory), new(PaymentItem), new(PaymentTry))

	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "postgres://sqlmcppd:rC_KcaIStkNyjO7rIRkVQTh77SFejZ7s@baasu.db.elephantsql.com:5432/sqlmcppd")

}
