package main

import (
	"coder/controller/app"
	"coder/packs/beego"
)

func init() {
	beego.BConfig.WebConfig.Session.SessionOn = true
}

func main() {
	beego.Router("/", &app.IndexController{}, "*:Index")
	beego.Router("/analyse", &app.IndexController{}, "*:Analyse")
	beego.Router("/logout", &app.LoginController{}, "get:Logout")
	beego.Router("/login", &app.LoginController{}, "get:Login")
	beego.Router("/trylogin", &app.LoginController{}, "post:DoLogin")
	beego.Run()
}