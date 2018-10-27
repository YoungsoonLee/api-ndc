package main

import (
	"os"
	"strconv"

	"github.com/YoungsoonLee/api-ndc/models"
	_ "github.com/YoungsoonLee/api-ndc/routers"
	"github.com/joho/godotenv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	models.RegisterDB()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		//log.Fatal("Error loading .env file")
		beego.Error("Error loading .env file")
	}

	// PORT
	PORT, _ := strconv.Atoi(os.Getenv("PORT"))
	if PORT == 0 {
		PORT, _ = strconv.Atoi(beego.AppConfig.String("httpport"))
	}
	beego.BConfig.Listen.HTTPPort = PORT

	// RUNMODE
	RUNMODE := os.Getenv("BEEGO_RUNMODE")
	if RUNMODE == "" {
		RUNMODE = beego.AppConfig.String("runmode")
	}
	beego.BConfig.RunMode = RUNMODE

	// beego.BConfig.Listen.HTTPSPort = 10443

	beego.Info("beego runmode: ", RUNMODE)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	}

	orm.RunSyncdb("default", false, true)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	beego.Run()
}
